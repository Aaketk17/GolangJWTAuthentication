package routes

import (
	controller "github.com/Aaketk17/GolangJWTAuthentication/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.SignUp)
	incomingRoutes.POST("users/login", controller.Login)
}
