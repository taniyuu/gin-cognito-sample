package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/taniyuu/gin-cognito-sample/domain/proxy"
)

// AuthzMiddleware アカウント認証操作を実行します
type AuthzMiddleware struct {
	ap proxy.AuthorizarProxy
}

// NewAuthzMiddleware AuthzMiddlewareを生成します
func NewAuthzMiddleware(ap proxy.AuthorizarProxy) *AuthzMiddleware {
	return &AuthzMiddleware{ap}
}

// Authorization アカウントを認証しコンテキストに設定します
func (am *AuthzMiddleware) Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.Replace(c.GetHeader("Authorization"), "Bearer ", "", 1)
		err := am.ap.ValidateJWT(token)
		if err != nil {
			am.errorResponse(c, err)
			c.Abort()
			return
		}
		c.Next()
	}
}
func (am *AuthzMiddleware) errorResponse(c *gin.Context, err error) {
	log.Default().Printf("%+v", err)
	// 適当なエラーレスポンス
	c.JSON(401, gin.H{
		"message": "unauthorized",
	})
}
