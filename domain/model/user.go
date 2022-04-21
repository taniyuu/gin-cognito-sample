package model

type User struct {
	Sub   string `json:"-"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type CreateReq struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
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
	// PreviousPassword string `json:"previous_password" validate:"required"`
	ProposedPassword string `json:"proposed_password" validate:"required"`
}

type ForgotPasswordReq struct {
	Email string `json:"email" validate:"required,email"`
}

type ConfirmForgotPasswordReq struct {
	Email    string `json:"email" validate:"required,email"`
	Code     string `json:"code" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ChangeProfileReq struct {
	Name string `json:"name" validate:"required"`
}

type SignoutReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type InviteReq struct {
	Email string `json:"email" validate:"required,email"`
}

type RespondToInvitationReq struct {
	Email            string `json:"email" validate:"required,email"`
	Name             string `json:"name" validate:"required"`
	Password         string `json:"password" validate:"required"`
	ConfirmationCode string `json:"confirmation_code" validate:"required"`
}

type Token struct {
	IDToken      string  `json:"id_token"`
	RefreshToken *string `json:"refresh_token,omitempty"`
}
