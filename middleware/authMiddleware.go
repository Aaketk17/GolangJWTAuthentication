package middleware

import (
	"net/http"

	helper "github.com/Aaketk17/GolangJWTAuthentication/helpers"
	"github.com/gin-gonic/gin"
)

func Authenticate(c *gin.Context) {
	clientToken := c.Request.Header.Get("token")
	if clientToken == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token is not valid"})
		c.Abort()
		return
	}

	claims, err := helper.ValidateToken(clientToken)
	if err != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		c.Abort()
		return
	}

	c.Set("email", claims.Email)
	c.Set("firstName", claims.FirstName)
	c.Set("lastName", claims.LastName)
	c.Set("uid", claims.Uid)
	c.Set("userType", claims.UserType)
	c.Next()
}
