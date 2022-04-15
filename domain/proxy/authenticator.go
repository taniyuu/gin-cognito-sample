package proxy

import (
	"context"

	"github.com/taniyuu/gin-cognito-sample/domain/model"
)

// AuthenticatorProxy 認証操作を抽象化します
type AuthenticatorProxy interface {
	Signup(ctx context.Context, req *model.CreateReq) (uuid string, err error)
	ConfirmAndSignin(ctx context.Context, req *model.ConfirmAndSigninReq) (*model.Token, error)
	Signin(ctx context.Context, req *model.SigninReq) (*model.Token, error)
	Refresh(ctx context.Context, req *model.RefreshReq) (*model.Token, error)
	ChangePassword(ctx context.Context, req *model.ChangePasswordReq) error
}
