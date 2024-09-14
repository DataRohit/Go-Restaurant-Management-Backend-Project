package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/datarohit/go-restaurant-management-backend-project/database"
	helper "github.com/datarohit/go-restaurant-management-backend-project/helpers"
	"github.com/datarohit/go-restaurant-management-backend-project/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

		totalCount, err := menuCollection.CountDocuments(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count menu items"})
			return
		}

		cursor, err := menuCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve menu items"})
			return
		}
		defer cursor.Close(ctx)

		var menus []models.Menu
		if err := cursor.All(ctx, &menus); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while parsing menu items"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"totalCount": totalCount,
			"menus":      menus,
		})
	}
}

func GetMenuByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		menuId := c.Param("menuId")
		if menuId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "menu_id parameter is required"})
			return
		}

		var menu models.Menu
		err := menuCollection.FindOne(ctx, bson.M{"menuId": menuId}).Decode(&menu)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the menu"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"menu": menu})
	}
}

func UpdateMenuByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		menuId := c.Param("menuId")
		if menuId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "menuId parameter is required"})
			return
		}

		var menu models.Menu
		if err := c.BindJSON(&menu); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}

		var updateObj primitive.D

		if menu.StartDate != nil && menu.EndDate != nil {
			if !helper.InTimeSpan(*menu.StartDate, *menu.EndDate, time.Now()) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Start date must be before end date"})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "startDate", Value: menu.StartDate}, bson.E{Key: "endDate", Value: menu.EndDate})
		}

		if menu.Name != "" {
			updateObj = append(updateObj, bson.E{Key: "name", Value: menu.Name})
		}

		if menu.Category != "" {
			updateObj = append(updateObj, bson.E{Key: "category", Value: menu.Category})
		}

		if len(updateObj) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
			return
		}

		menu.UpdatedAt = time.Now()
		updateObj = append(updateObj, bson.E{Key: "updated_at", Value: menu.UpdatedAt})

		filter := bson.M{"menuId": menuId}
		opts := options.Update().SetUpsert(true)

		_, err := menuCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update menu"})
			return
		}

		var updatedMenu models.Menu
		if err := menuCollection.FindOne(ctx, filter).Decode(&updatedMenu); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated menu"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Menu updated successfully", "menu": updatedMenu})
	}
}
