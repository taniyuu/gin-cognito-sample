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

type RefreshReq struct {
	Sub          string `json:"sub"`
	RefreshToken string `json:"refresh_token"`
}

type Token struct {
	IDToken      string  `json:"id_token"`
	RefreshToken *string `json:"refresh_token,omitempty"`
}
