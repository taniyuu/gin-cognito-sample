package usecase

import (
	"context"

	"github.com/taniyuu/gin-cognito-sample/application/viewmodel"
	"github.com/taniyuu/gin-cognito-sample/domain/proxy"
)

// UserUsecase アカウントに対する操作を抽象化します
type UserUsecase interface {
	Create(ctx context.Context, req *viewmodel.CreateReq) error
	Confirm(ctx context.Context, req *viewmodel.ConfirmReq) error
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
	_, err := tu.ap.Signup(ctx, &req.CreateReq)
	// uuidを返すので、必要であれば利用
	return err
}

// Create アカウント新規作成
func (tu *userUsecase) Confirm(ctx context.Context, req *viewmodel.ConfirmReq) error {
	_, err := tu.ap.Confirm(ctx, &req.ConfirmReq)
	return err
}
