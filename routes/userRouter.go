package routes

import (
	controller "github.com/Aaketk17/GolangJWTAuthentication/controllers"
	middleware "github.com/Aaketk17/GolangJWTAuthentication/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Athenticate())
	incomingRoutes.GET("users", controller.GetUsers)
	incomingRoutes.GET("users/:user_id", controller.GetUser)
}