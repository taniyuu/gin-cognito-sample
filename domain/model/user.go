package model

type UserInfo struct {
	Email string
	Name  string
	Id    string
}

type CreateReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ConfirmAndSigninReq struct {
	Email            string `json:"email" validate:"required,email"`
	ConfirmationCode string `json:"confirmation_code" validate:"required"`
	Password         string `json:"password" validate:"required"`
}

type SigninReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
