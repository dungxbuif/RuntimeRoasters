package errs

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
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

// StatusCode returns the HTTP status code for a given error
func StatusCode(err error) int {
	if s, ok := HTTPStatusMap[err]; ok {
		return s
	}
	return http.StatusInternalServerError
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

		p := Problem{
			Type:     "about:blank",
			Title:    http.StatusText(status),
			Status:   status,
			Detail:   err.Error(),
			Instance: c.Request.URL.Path,
			TraceID:  c.GetHeader("X-Trace-ID"),
		}

		c.Render(status, problemRenderer{p})
	}
}
