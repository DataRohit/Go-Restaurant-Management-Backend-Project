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
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func CreateFood() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var food models.Food
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}

		if err := validate.Struct(food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var menu models.Menu
		err := menuCollection.FindOne(ctx, bson.M{"menuId": food.MenuID}).Decode(&menu)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching menu"})
			return
		}

		food.CreatedAt = time.Now().UTC()
		food.UpdatedAt = time.Now().UTC()
		food.ID = primitive.NewObjectID()
		food.FoodID = food.ID.Hex()

		roundedPrice := helper.ToFixed(*food.Price, 2)
		food.Price = &roundedPrice

		_, err = foodCollection.InsertOne(ctx, food)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create food item"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Food item created successfully", "food": food})
	}
}

func GetAllFoodItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		recordPerPage, page := helper.GetPaginationParams(c)

		skip := (page - 1) * recordPerPage
		foodItems, totalCount, err := helper.GetPaginatedFoodItems(ctx, skip, recordPerPage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"totalCount": totalCount,
			"foodItems":  foodItems,
		})
	}
}
