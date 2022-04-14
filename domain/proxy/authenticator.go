package proxy

import (
	"context"

	"github.com/taniyuu/gin-cognito-sample/domain/model"
)

// AuthenticatorProxy 認証操作を抽象化します
type AuthenticatorProxy interface {
	Signup(ctx context.Context, info *model.UserInfo) (string, error)
}
