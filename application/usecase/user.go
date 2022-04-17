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
	Signin(ctx context.Context, req *viewmodel.SigninReq) (*viewmodel.SigninResp, error)
	Refresh(ctx context.Context, req *viewmodel.RefreshReq) (*viewmodel.SigninResp, error)
	ChangePassword(ctx context.Context, req *viewmodel.ChangePasswordReq) error
	GetProfile(ctx context.Context, req *viewmodel.GetProfileReq) (*viewmodel.User, error)
	Signout(ctx context.Context, req *viewmodel.SignoutReq) error
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

// Signin アカウント確認を行います（ログインも試行する、MFAが設定された認証プールには適用できないので注意）
func (tu *userUsecase) Signin(ctx context.Context, req *viewmodel.SigninReq) (*viewmodel.SigninResp, error) {
	token, err := tu.ap.Signin(ctx, &req.SigninReq)
	if err != nil {
		return nil, err
	}
	resp := new(viewmodel.SigninResp)
	resp.Token = *token
	return resp, nil
}

// Refresh トークンリフレッシュを行います
func (tu *userUsecase) Refresh(ctx context.Context, req *viewmodel.RefreshReq) (*viewmodel.SigninResp, error) {
	token, err := tu.ap.Refresh(ctx, &req.RefreshReq)
	if err != nil {
		return nil, err
	}
	resp := new(viewmodel.SigninResp)
	resp.Token = *token
	return resp, nil
}

// ChangePassword パスワード変更を行います
func (tu *userUsecase) ChangePassword(ctx context.Context, req *viewmodel.ChangePasswordReq) error {
	return tu.ap.ChangePassword(ctx, &req.ChangePasswordReq)
}

// GetProfile アカウント情報を取得します
func (tu *userUsecase) GetProfile(ctx context.Context, req *viewmodel.GetProfileReq) (*viewmodel.User, error) {
	user, err := tu.ap.GetProfile(ctx, &req.GetProfileReq)
	if err != nil {
		return nil, err
	}
	resp := new(viewmodel.User)
	resp.User = *user
	return resp, nil
}

// Signout ログアウトを行います
func (tu *userUsecase) Signout(ctx context.Context, req *viewmodel.SignoutReq) error {
	return tu.ap.Signout(ctx, &req.SignoutReq)
}
