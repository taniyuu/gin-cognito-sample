package middleware

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/taniyuu/gin-cognito-sample/domain/proxy"
)

const subContextKey string = "sub"

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
		token := c.GetHeader("Authorization")
		sub, err := am.ap.ValidateJWT(token)
		if err != nil {
			am.errorResponse(c, err)
			c.Abort()
			return
		}
		// ginコンテキストにsubを入れる
		c.Set(subContextKey, sub)
		c.Next()
	}
}

func GetSub(c *gin.Context) (string, error) {
	v := c.GetString(subContextKey)
	if v == "" {
		return v, errors.WithStack(fmt.Errorf("token not found"))
	}
	return v, nil
}

func (am *AuthzMiddleware) errorResponse(c *gin.Context, err error) {
	log.Default().Printf("%+v", err)
	// 適当なエラーレスポンス
	c.JSON(401, gin.H{
		"message": "unauthorized",
	})
}
