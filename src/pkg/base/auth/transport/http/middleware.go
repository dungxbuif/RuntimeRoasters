package authhttp

import (
	"net/http"
	"strings"

	"RuntimeRoasters/pkg/base/auth/provider"
	"RuntimeRoasters/pkg/base/auth/token"
	"RuntimeRoasters/pkg/errs"
	"RuntimeRoasters/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MiddlewareOptions struct {
	PublicRoutes     []string
	BlacklistChecker token.BlacklistChecker
}

type MiddlewareOption func(*MiddlewareOptions)

func WithPublicRoutes(routes ...string) MiddlewareOption {
	return func(o *MiddlewareOptions) {
		o.PublicRoutes = append(o.PublicRoutes, routes...)
	}
}

func WithBlacklistChecker(checker token.BlacklistChecker) MiddlewareOption {
	return func(o *MiddlewareOptions) {
		o.BlacklistChecker = checker
	}
}

func GinMiddleware(keyProvider provider.KeyProvider, expectedIssuer string, opts ...MiddlewareOption) gin.HandlerFunc {
	options := &MiddlewareOptions{
		PublicRoutes: []string{"/health/live", "/health/ready"},
	}
	for _, opt := range opts {
		opt(options)
	}

	public := make(map[string]struct{}, len(options.PublicRoutes))
	for _, route := range options.PublicRoutes {
		public[route] = struct{}{}
	}

	return func(c *gin.Context) {
		if _, ok := public[c.FullPath()]; ok {
			c.Next()
			return
		}
		if _, ok := public[c.Request.URL.Path]; ok {
			c.Next()
			return
		}

		raw, ok := bearerToken(c.GetHeader(token.HeaderAuthorization))
		if !ok {
			_ = c.Error(errs.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, err := token.VerifyAndParseJWT(raw, keyProvider, expectedIssuer)
		if err != nil {
			logger.FromContext(c.Request.Context()).Warn("JWT verification failed", zap.Error(err))
			_ = c.Error(errs.ErrUnauthorized)
			c.Abort()
			return
		}

		if options.BlacklistChecker != nil {
			revoked, err := options.BlacklistChecker.IsRevoked(c.Request.Context(), claims.JTI)
			if err != nil {
				logger.FromContext(c.Request.Context()).Error("blacklist check failed", zap.Error(err), zap.String("jti", claims.JTI))
				_ = c.Error(errs.ErrInternal)
				c.Abort()
				return
			}
			if revoked {
				_ = c.Error(errs.ErrUnauthorized)
				c.Abort()
				return
			}
		}

		ctx := token.SetIdentityInContext(c.Request.Context(), *claims)
		token.RecordTracingData(ctx, claims.Subject)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func bearerToken(header string) (string, bool) {
	parts := strings.Split(header, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != token.BearerPrefix || strings.TrimSpace(parts[1]) == "" {
		return "", false
	}
	return parts[1], true
}

func AbortUnauthorized(c *gin.Context) {
	c.AbortWithStatus(http.StatusUnauthorized)
}
