package handler

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/taniyuu/gin-cognito-sample/application/usecase"
	"github.com/taniyuu/gin-cognito-sample/application/viewmodel"
	"gopkg.in/go-playground/validator.v9"
)

type UserHandler struct {
	tu usecase.UserUsecase
	v  *validator.Validate
}

func NewUserHandler(tu usecase.UserUsecase) *UserHandler {
	return &UserHandler{tu, validator.New()}
}

func (h *UserHandler) Create(c *gin.Context) {
	req := new(viewmodel.CreateReq)
	if err := c.ShouldBindJSON(req); err != nil {
		h.errorResponse(c, err)
		return
	}
	if err := h.v.Struct(req); err != nil {
		h.errorResponse(c, err)
		return
	}

	err := h.tu.Create(c.Request.Context(), req)
	if err != nil {
		h.errorResponse(c, err)
	} else {
		c.Status(200)
	}
}

func (h *UserHandler) errorResponse(c *gin.Context, err error) {
	log.Default().Printf("%+v", err)
	// 適当なエラーレスポンス
	c.JSON(500, gin.H{
		"message": "server error",
	})
}
