package viewmodel

import "github.com/taniyuu/gin-cognito-sample/domain/model"

type CreateReq struct {
	model.CreateReq
}

type ConfirmReq struct {
	model.ConfirmAndSigninReq
}

type SigninReq struct {
	model.SigninReq
}

type RefreshReq struct {
	model.RefreshReq
}

type ChangePasswordReq struct {
	model.ChangePasswordReq
}

type GetProfileReq struct {
	model.GetProfileReq
}
type SigninResp struct {
	model.Token
}

type User struct {
	model.User
}
