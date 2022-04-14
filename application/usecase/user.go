package usecase

import (
	"context"

	"github.com/taniyuu/gin-cognito-sample/domain/model"
	"github.com/taniyuu/gin-cognito-sample/domain/proxy"
)

// UserUsecase アカウントに対する操作を抽象化します
type UserUsecase interface {
	Create(ctx context.Context, email string) error
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
func (tu *userUsecase) Create(ctx context.Context, email string) error {
	_, err := tu.ap.Signup(ctx, &model.UserInfo{})
	return err
}
