package main

import (
	"os"

	routes "github.com/Aaketk17/GolangJWTAuthentication/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "8080"
	}

	router := gin.Default()

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.Run(":" + PORT)

}
