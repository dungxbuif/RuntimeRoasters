package errs

import (
	"errors"
	"net/http"

	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const problemContentType = "application/problem+json"

// Problem struct (RFC 9457)
type Problem struct {
	Type     string         `json:"type"`
	Title    string         `json:"title"`
	Status   int            `json:"status"`
	Detail   string         `json:"detail,omitempty"`
	Instance string         `json:"instance,omitempty"`
	TraceID  string         `json:"trace_id,omitempty"`
	Errors   []SubProblem   `json:"errors,omitempty"`
}

// SubProblem represents a single entry in the errors extension (RFC 9457 §3.1)
type SubProblem struct {
	Detail  string `json:"detail"`
	Pointer string `json:"pointer,omitempty"` // JSON Pointer (RFC 6901) to the offending field
}

// Sentinel errors
var (
	ErrNotFound     = errors.New("not_found")
	ErrConflict     = errors.New("conflict")
	ErrUnauthorized = errors.New("unauthorized")
	ErrValidation   = errors.New("validation_failed")
	ErrForbidden    = errors.New("forbidden")
	ErrInternal     = errors.New("internal_error")
)

// HTTPStatusMap mapping sentinel errors to HTTP status codes
var HTTPStatusMap = map[error]int{
	ErrNotFound:     http.StatusNotFound,
	ErrConflict:     http.StatusConflict,
	ErrUnauthorized: http.StatusUnauthorized,
	ErrValidation:   http.StatusBadRequest,
	ErrForbidden:    http.StatusForbidden,
	ErrInternal:     http.StatusInternalServerError,
}

// GRPCStatusMap mapping sentinel errors to gRPC codes
var GRPCStatusMap = map[error]codes.Code{
	ErrNotFound:     codes.NotFound,
	ErrConflict:     codes.Aborted,
	ErrUnauthorized: codes.Unauthenticated,
	ErrValidation:   codes.InvalidArgument,
	ErrForbidden:    codes.PermissionDenied,
	ErrInternal:     codes.Internal,
}

// StatusCode returns the HTTP status code for a given error
func StatusCode(err error) int {
	// Check if it's a wrapped error
	for sentinel, code := range HTTPStatusMap {
		if errors.Is(err, sentinel) {
			return code
		}
	}
	return http.StatusInternalServerError
}

// GRPCCode returns the gRPC code for a given error
func GRPCCode(err error) codes.Code {
	for sentinel, code := range GRPCStatusMap {
		if errors.Is(err, sentinel) {
			return code
		}
	}
	return codes.Internal
}

// ToGRPCError converts a sentinel error to a gRPC status error
func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}
	return status.Error(GRPCCode(err), err.Error())
}

// New builds a Problem with type defaulting to "about:blank" per RFC 9457
func New(status int, detail string) Problem {
	return Problem{
		Type:   "about:blank",
		Title:  http.StatusText(status),
		Status: status,
		Detail: detail,
	}
}

// GinErrorHandler middleware - converts c.Error(err) to RFC 9457 JSON response
func GinErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		status := StatusCode(err)
		
		// Extract TraceID from OpenTelemetry
		traceID := trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String()

		log := logger.FromContext(c.Request.Context())
		if status >= 500 {
			log.Error("internal server error", zap.Error(err), zap.String("path", c.Request.URL.Path))
		} else {
			log.Warn("client error", zap.Error(err), zap.Int("status", status), zap.String("path", c.Request.URL.Path))
		}

		p := Problem{
			Type:     "about:blank",
			Title:    http.StatusText(status),
			Status:   status,
			Detail:   err.Error(),
			Instance: c.Request.URL.Path,
			TraceID:  traceID,
		}

		c.Render(status, problemRenderer{p})
	}
}
