package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Aaketk17/GolangJWTAuthentication/models"

	"github.com/Aaketk17/GolangJWTAuthentication/database"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	helper "github.com/Aaketk17/GolangJWTAuthentication/helpers"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func HashPassword() {}

func VerifyPassword() {}

func SignUp(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User
	defer cancel()

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := validate.Struct(user) //!
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailCount, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if emailCount > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "This email adddress already taken"})
	}

	phoneCount, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	// defer cancel()  //!
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if phoneCount > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "This Phone already taken"})
	}

}

func Login() {}

func GetUsers() {}

func GetUser(c *gin.Context) {
	userId := c.Param("user_id")

	err := helper.MatchUserTypeTpUid(c, userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"userId": userId}).Decode(&user)
	defer cancel()
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record does not exists"})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)

}
