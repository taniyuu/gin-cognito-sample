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
	ChangePassword(ctx context.Context, email string, req *viewmodel.ChangePasswordReq) error
	ForgotPassword(ctx context.Context, req *viewmodel.ForgotPasswordReq) error
	ConfirmForgotPassword(ctx context.Context, req *viewmodel.ConfirmForgotPasswordReq) error
	GetProfile(ctx context.Context, email string) (*viewmodel.User, error)
	ChangeProfile(ctx context.Context, email string, req *viewmodel.ChangeProfileReq) error
	Signout(ctx context.Context, req *viewmodel.SignoutReq) error
	Invite(ctx context.Context, req *viewmodel.InviteReq) (*viewmodel.InviteResp, error)
	RespondToInvitation(ctx context.Context, req *viewmodel.RespondToInvitationReq) (*viewmodel.SigninResp, error)
	GetUserForAdmin(ctx context.Context, id string) (*viewmodel.User, error)
}

// アカウントに対する操作を提供します
type userUsecase struct {
	ap proxy.UserProxy
}

// NewUserUsecase UserUsecaseを生成します
func NewUserUsecase(
	ap proxy.UserProxy,
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
func (tu *userUsecase) ChangePassword(ctx context.Context, email string, req *viewmodel.ChangePasswordReq) error {
	return tu.ap.ChangePassword(ctx, email, &req.ChangePasswordReq)
}

// ForgotPassword パスワード忘れ
func (tu *userUsecase) ForgotPassword(ctx context.Context, req *viewmodel.ForgotPasswordReq) error {
	return tu.ap.ForgotPassword(ctx, &req.ForgotPasswordReq)
}

// ConfirmForgotPassword パスワード忘れ確認
func (tu *userUsecase) ConfirmForgotPassword(ctx context.Context, req *viewmodel.ConfirmForgotPasswordReq) error {
	return tu.ap.ConfirmForgotPassword(ctx, &req.ConfirmForgotPasswordReq)
}

// GetProfile アカウント情報を取得します
func (tu *userUsecase) GetProfile(ctx context.Context, email string) (*viewmodel.User, error) {
	user, err := tu.ap.GetProfile(ctx, email)
	if err != nil {
		return nil, err
	}
	resp := new(viewmodel.User)
	resp.User = *user
	return resp, nil
}

// ChangeProfile アカウント情報を変更します
func (tu *userUsecase) ChangeProfile(ctx context.Context, email string, req *viewmodel.ChangeProfileReq) error {
	return tu.ap.ChangeProfile(ctx, email, &req.ChangeProfileReq)
}

// Signout ログアウトを行います
func (tu *userUsecase) Signout(ctx context.Context, req *viewmodel.SignoutReq) error {
	return tu.ap.Signout(ctx, &req.SignoutReq)
}

// Invite 招待を行います
func (tu *userUsecase) Invite(ctx context.Context, req *viewmodel.InviteReq) (*viewmodel.InviteResp, error) {
	sub, err := tu.ap.Invite(ctx, &req.InviteReq)
	if err != nil {
		return nil, err
	}
	return &viewmodel.InviteResp{Sub: sub}, nil
}

// RespondToInvitation 招待応答を行います
func (tu *userUsecase) RespondToInvitation(ctx context.Context, req *viewmodel.RespondToInvitationReq) (*viewmodel.SigninResp, error) {
	token, err := tu.ap.RespondToInvitation(ctx, &req.RespondToInvitationReq)
	if err != nil {
		return nil, err
	}
	resp := new(viewmodel.SigninResp)
	resp.Token = *token
	return resp, nil
}

// GetUserForAdmin ユーザ取得を行います
func (tu *userUsecase) GetUserForAdmin(ctx context.Context, id string) (*viewmodel.User, error) {
	user, err := tu.ap.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := new(viewmodel.User)
	resp.User = *user
	return resp, nil
}
