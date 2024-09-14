package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/datarohit/go-restaurant-management-backend-project/database"
	"github.com/datarohit/go-restaurant-management-backend-project/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var menuCollection *mongo.Collection = database.OpenCollection(database.Client, "menu")

func CreateMenu() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var menu models.Menu
		if err := c.ShouldBindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		if err := validate.Struct(menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		currentTime := time.Now().UTC()
		menu.CreatedAt = currentTime
		menu.UpdatedAt = currentTime
		menu.ID = primitive.NewObjectID()
		menu.MenuID = menu.ID.Hex()

		if _, err := menuCollection.InsertOne(ctx, menu); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create menu item"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Menu item created successfully", "menu": menu})
	}
}

func GetAllMenus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := menuCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve menu items"})
			return
		}
		defer cursor.Close(ctx)

		var menus []bson.M
		if err := cursor.All(ctx, &menus); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while parsing menu items"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"menus": menus})
	}
}
