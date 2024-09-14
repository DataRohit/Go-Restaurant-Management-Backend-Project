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

func GetFoodByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		foodID := c.Param("foodId")
		if foodID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "foodId parameter is required"})
			return
		}

		var food models.Food
		err := foodCollection.FindOne(ctx, bson.M{"foodId": foodID}).Decode(&food)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Food item not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while fetching the food item"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"food": food})
	}
}

func UpdateFoodByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		foodID := c.Param("foodId")
		if foodID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "foodId parameter is required"})
			return
		}

		var food models.Food
		if err := c.BindJSON(&food); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}

		updateObj := bson.D{}

		if food.Name != nil {
			updateObj = append(updateObj, bson.E{Key: "name", Value: food.Name})
		}

		if food.Price != nil {
			roundedPrice := helper.ToFixed(*food.Price, 2)
			updateObj = append(updateObj, bson.E{Key: "price", Value: roundedPrice})
		}

		if food.FoodImage != nil {
			updateObj = append(updateObj, bson.E{Key: "foodImage", Value: food.FoodImage})
		}

		if food.MenuID != nil {
			var menu models.Menu
			err := menuCollection.FindOne(ctx, bson.M{"menuId": food.MenuID}).Decode(&menu)
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Menu not found"})
				return
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching menu"})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "menuId", Value: food.MenuID})
		}

		if len(updateObj) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
			return
		}

		food.UpdatedAt = time.Now().UTC()
		updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: food.UpdatedAt})

		filter := bson.M{"foodId": foodID}
		opts := options.Update().SetUpsert(false)

		result, err := foodCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update food item"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Food item not found"})
			return
		}

		var updatedFood models.Food
		if err := foodCollection.FindOne(ctx, filter).Decode(&updatedFood); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated food item"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Food item updated successfully", "food": updatedFood})
	}
}
