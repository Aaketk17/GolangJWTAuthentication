package main

import (
	"os"

	routes "github.com/Aaketk17/GolangJWTAuthentication/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	router := gin.Default()

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.Run(port)

}
