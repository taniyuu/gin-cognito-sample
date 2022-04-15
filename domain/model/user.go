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
