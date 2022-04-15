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

type SigninResp struct {
	model.Token
}
