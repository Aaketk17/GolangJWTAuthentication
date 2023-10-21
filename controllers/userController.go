package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

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

func HashPassword(userPwd string) string {
	hashValue, err := bcrypt.GenerateFromPassword([]byte(userPwd), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(hashValue)
}

func VerifyPassword(userPwd string, proviedPwd string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(proviedPwd), []byte(userPwd))
	check := true
	msg := ""

	if err != nil {
		msg = "email of password is incorrect"
		check = false
	}

	return check, msg
}

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
		return
	}

	phoneCount, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	// defer cancel()  //!
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if phoneCount > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "This Phone already taken"})
		return
	}

	userPwd := HashPassword(*user.Password)
	user.Password = &userPwd

	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	user.ID = primitive.NewObjectID()
	user.UserID = user.ID.Hex()
	token, refreshToken, _ := helper.GenerateTokens(*user.Email, *user.FirstName, *user.LastName, *user.UserType, user.UserID)
	user.Token = &token
	user.RefreshToken = &refreshToken

	insertionNumber, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in creating user"})
		return
	}
	c.JSON(http.StatusOK, insertionNumber)

}

func Login(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User
	var foundUser models.User

	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	findErr := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)

	if findErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": findErr.Error()})
		return
	}

	pwdValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

	if !pwdValid {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	if foundUser.Email == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
	} //!

	token, refreshToken, _ := helper.GenerateTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, *foundUser.UserType, foundUser.UserID)
	helper.UpdateAllTokens(token, refreshToken, foundUser.UserID)

	findErrTwo := userCollection.FindOne(ctx, bson.M{"userid": foundUser.UserID}).Decode(&foundUser) //!
	if findErrTwo != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": findErr.Error()})
		return
	}
	c.JSON(http.StatusOK, foundUser)
}

func GetUsers(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err := helper.CheckUserType(c, "ADMIN")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}
	page, err1 := strconv.Atoi(c.Query("page"))
	if err1 != nil || page < 1 {
		page = 1
	}

	startIndex := (page - 1) * recordPerPage
	startIndex, err = strconv.Atoi(c.Query("startIndex"))

	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
		{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
		{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "user_items", Value: bson.D{{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}}

	result, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})

	defer cancel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
	}
	var allusers []bson.M
	if err = result.All(ctx, &allusers); err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, allusers[0])
}

func GetUser(c *gin.Context) {
	userId := c.Param("user_id")

	err := helper.MatchUserTypeTpUid(c, userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"userid": userId}).Decode(&user)
	defer cancel()
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Record does not exists"})
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)

}
