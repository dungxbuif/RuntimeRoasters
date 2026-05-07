package auth

import (
	"strings"

	"github.com/dungxbuif/RuntimeRoasters/pkg/errs"
	"github.com/gin-gonic/gin"
)

/* Dùng JWTs để auth
- JWT truyền từ bên ngoài vào để tách biệt logic auth và logic business
- Token là dạng chữ ký xác thực, UseCase sẽ không nhìn thấy Token mà chỉ thấy ID
*/

func extractTokenFromHeader(ctx *gin.Context) (string, bool) {
	authHeader := ctx.GetHeader(HeaderAuthorization)
	if len(authHeader) == 0 {
		ctx.Error(errs.ErrUnauthorized).SetMeta("missing authorization header")
		ctx.Abort()
		return "", false
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != BearerPrefix {
		ctx.Error(errs.ErrUnauthorized).SetMeta("invalid authorization header format")
		ctx.Abort()
		return "", false
	}
	return parts[1], true
}

func GinMiddleware(keyProvider KeyProvider, expectedIssuer string) gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtToken, ok := extractTokenFromHeader(c)
		if !ok {
			return
		}

		token, err := VerifyAndParseJWT(jwtToken, keyProvider, expectedIssuer)
		if err != nil {
			c.Error(errs.ErrUnauthorized).SetMeta(err.Error())
			c.Abort()
			return
		}

		ctx := SetIdentityInContext(c.Request.Context(), *token)
		c.Request = c.Request.WithContext(ctx)

		RecordTracingData(ctx, token.Subject)

		c.Next()
	}
}