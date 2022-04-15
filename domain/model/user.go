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
	Sub          string `json:"sub" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ChangePasswordReq struct {
	AccessToken      string `json:"access_token" validate:"required"`
	PreviousPassword string `json:"previous_password" validate:"required"`
	ProposedPassword string `json:"proposed_password" validate:"required"`
}

type Token struct {
	IDToken      string  `json:"id_token"`
	AccessToken  string  `json:"access_token"`
	RefreshToken *string `json:"refresh_token,omitempty"`
}
