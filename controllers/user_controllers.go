package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/datarohit/go-restaurant-management-backend-project/database"
	helper "github.com/datarohit/go-restaurant-management-backend-project/helpers"
	"github.com/datarohit/go-restaurant-management-backend-project/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validate.Struct(user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		emailExists, err := helper.RecordExists(ctx, userCollection, "email", user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking email existence"})
			return
		}

		if emailExists {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}

		phoneExists, err := helper.RecordExists(ctx, userCollection, "phone", user.Phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking phone existence"})
			return
		}

		if phoneExists {
			c.JSON(http.StatusConflict, gin.H{"error": "Phone number already exists"})
			return
		}

		user.Password = helper.HashPassword(user.Password)
		user.CreatedAt = time.Now()
		user.UpdatedAt = user.CreatedAt
		user.ID = primitive.NewObjectID()
		user.UserID = user.ID.Hex()

		accessToken, refreshToken, err := helper.GenerateAllTokens(*user.Email, *user.FirstName, *user.LastName, user.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating tokens"})
			return
		}

		user.AccessToken = &accessToken
		user.RefreshToken = &refreshToken

		if _, err := userCollection.InsertOne(ctx, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	}
}
