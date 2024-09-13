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
	"go.mongodb.org/mongo-driver/bson"
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

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		passwordIsValid, msg := helper.VerifyPassword(*foundUser.Password, *user.Password)
		if !passwordIsValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": msg})
			return
		}

		var accessToken, refreshToken string
		var generateNewTokens bool

		if err := helper.ValidateToken(*foundUser.AccessToken); err == nil {
			accessToken = *foundUser.AccessToken
		} else {
			if err := helper.ValidateToken(*foundUser.RefreshToken); err == nil {
				newAccessToken, err := helper.GenerateToken(*foundUser.Email, foundUser.UserID, 24*time.Hour)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating new access token"})
					return
				}
				accessToken = newAccessToken
			} else {
				newAccessToken, newRefreshToken, err := helper.GenerateAllTokens(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.UserID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating tokens"})
					return
				}
				accessToken = newAccessToken
				refreshToken = newRefreshToken
				generateNewTokens = true
			}
		}

		if generateNewTokens {
			err := helper.UpdateAllTokens(accessToken, refreshToken, foundUser.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating tokens"})
				return
			}
		} else {
			err := helper.UpdateAccessToken(accessToken, foundUser.UserID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating access token"})
				return
			}
			refreshToken = *foundUser.RefreshToken
		}

		c.JSON(http.StatusOK, gin.H{
			"user":         foundUser,
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		})
	}
}
