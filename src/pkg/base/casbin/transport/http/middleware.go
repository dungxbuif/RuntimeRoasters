package casbinhttp

import (
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/casbin"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/dungxbuif/RuntimeRoasters/pkg/errs"
	"github.com/gin-gonic/gin"
)

// GinMiddleware creates a Gin middleware for Casbin authorization
func GinMiddleware(engine casbin.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get identity from context
		claims, ok := identity.FromContext(c.Request.Context())
		if !ok {
			c.Error(errs.ErrUnauthorized)
			c.Abort()
			return
		}

		// 2. Map HTTP to Action
		action := casbin.MapHTTPMethodToAction(c.Request.Method)
		
		// 3. Object mapping
		// For simplicity, we use the URL path as the object
		// or map it to the gRPC method equivalent if possible.
		obj := c.Request.URL.Path

		// 4. Enforce
		allowed, err := engine.Enforce(claims.Role, obj, action)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		if !allowed {
			c.Error(errs.ErrForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}
