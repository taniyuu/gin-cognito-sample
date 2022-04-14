package handler

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/taniyuu/gin-cognito-sample/application/usecase"
)

type UserHandler struct {
	tu usecase.UserUsecase
}

func NewUserHandler(tu usecase.UserUsecase) *UserHandler {
	return &UserHandler{tu}
}

func (h *UserHandler) Create(c *gin.Context) {
	err := h.tu.Create(c.Request.Context(), "mail")
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
