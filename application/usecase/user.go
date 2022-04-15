package usecase

import (
	"context"

	"github.com/taniyuu/gin-cognito-sample/application/viewmodel"
	"github.com/taniyuu/gin-cognito-sample/domain/proxy"
)

// UserUsecase アカウントに対する操作を抽象化します
type UserUsecase interface {
	Create(ctx context.Context, req *viewmodel.CreateReq) error
	Confirm(ctx context.Context, req *viewmodel.ConfirmReq) (*viewmodel.SigninResp, error)
}

// アカウントに対する操作を提供します
type userUsecase struct {
	ap proxy.AuthenticatorProxy
}

// NewUserUsecase UserUsecaseを生成します
func NewUserUsecase(
	ap proxy.AuthenticatorProxy,
) UserUsecase {
	return &userUsecase{ap}
}

// Create アカウント新規作成
func (tu *userUsecase) Create(ctx context.Context, req *viewmodel.CreateReq) error {
	// uuidを返すので、利用可能
	_, err := tu.ap.Signup(ctx, &req.CreateReq)
	return err
}

// Confirm アカウント確認を行います（ログインも試行する、MFAが設定された認証プールには適用できないので注意）
func (tu *userUsecase) Confirm(ctx context.Context, req *viewmodel.ConfirmReq) (*viewmodel.SigninResp, error) {
	token, err := tu.ap.ConfirmAndSignin(ctx, &req.ConfirmAndSigninReq)
	if err != nil {
		return nil, err
	}
	resp := new(viewmodel.SigninResp)
	resp.Token = *token
	return resp, nil
}
