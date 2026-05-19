package casbinhttp

import (
	"net/http"

	"github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/dungxbuif/RuntimeRoasters/pkg/errs"
	"github.com/gin-gonic/gin"
)

type Enforcer interface {
	Enforce(rvals ...interface{}) (bool, error)
}

type RouteObjectFunc func(*gin.Context) string

type MiddlewareOptions struct {
	PublicRoutes []string
	RouteObject  RouteObjectFunc
}

type MiddlewareOption func(*MiddlewareOptions)

func WithPublicRoutes(routes ...string) MiddlewareOption {
	return func(o *MiddlewareOptions) {
		o.PublicRoutes = append(o.PublicRoutes, routes...)
	}
}

func WithRouteObject(fn RouteObjectFunc) MiddlewareOption {
	return func(o *MiddlewareOptions) {
		o.RouteObject = fn
	}
}

func GinMiddleware(engine Enforcer, opts ...MiddlewareOption) gin.HandlerFunc {
	options := &MiddlewareOptions{
		PublicRoutes: []string{"/health/live", "/health/ready"},
		RouteObject: func(c *gin.Context) string {
			if fullPath := c.FullPath(); fullPath != "" {
				return fullPath
			}
			return c.Request.URL.Path
		},
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

		claims, ok := identity.FromContext(c.Request.Context())
		if !ok {
			_ = c.Error(errs.ErrUnauthorized)
			c.Abort()
			return
		}

		object := options.RouteObject(c)
		action := casbin.MapHTTPMethodToAction(c.Request.Method)
		allowed, err := engine.Enforce(claims.Role, object, action)
		if err != nil {
			_ = c.Error(errs.ErrInternal)
			c.Abort()
			return
		}
		if !allowed {
			_ = c.Error(errs.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}

func AbortForbidden(c *gin.Context) {
	c.AbortWithStatus(http.StatusForbidden)
}
