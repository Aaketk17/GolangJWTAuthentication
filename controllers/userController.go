package controllers

import (
	database "github.com/Aaketk17/GolangJWTAuthentication/database"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func HashPassword() {}

func VerifyPassword() {}

func SignUp() {}

func Login() {}

func GetUsers() {}

func GetUser() {}
