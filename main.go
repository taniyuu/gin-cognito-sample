package main

import (
	"log"
	"net/http"
	"os"

	"github.com/taniyuu/gin-cognito-sample/application/usecase"
	awsWrapper "github.com/taniyuu/gin-cognito-sample/infrastructure/aws"
	"github.com/taniyuu/gin-cognito-sample/interface/handler"
	"github.com/taniyuu/gin-cognito-sample/interface/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cp := awsWrapper.NewCognitoProxy(
		os.Getenv("COGNITO_POOL_ID"), os.Getenv("COGNITO_CLIENT_ID"), os.Getenv("COGNITO_CLIENT_SECRET"))
	ap := awsWrapper.NewCognitoAuthorizar(os.Getenv("COGNITO_REGION"), os.Getenv("COGNITO_POOL_ID"), os.Getenv("COGNITO_CLIENT_ID"))
	uu := usecase.NewUserUsecase(cp)
	uh, am := handler.NewUserHandler(uu), middleware.NewAuthzMiddleware(ap)

	engine := gin.Default()
	// no authenticate endpoint
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "hello world",
		})
	})
	engine.POST("/signup", uh.Create)
	engine.POST("/confirm-signup", uh.Confirm)
	engine.POST("/signin", uh.Signin)
	authz := engine.Group("/", am.Authorization())
	{
		authz.GET("/check", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "hello world",
			})
		})
	}
	engine.Run(":3000")
}
