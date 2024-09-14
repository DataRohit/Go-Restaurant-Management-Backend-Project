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

type OrderItemPack struct {
	TableID    string `json:"tableId" binding:"required"`
	OrderItems []struct {
		FoodID   string `json:"foodId" binding:"required"`
		Quantity string `json:"quantity" binding:"required,oneof=S M L"`
	} `json:"orderItems" binding:"required"`
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func CreateOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItemPack OrderItemPack

		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		order := models.Order{
			OrderDate: time.Now().UTC(),
			TableID:   &orderItemPack.TableID,
		}

		orderId, err := helper.OrderItemOrderCreator(order, ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
			return
		}

		var orderItemsToBeInserted []interface{}
		var createdOrderItems []models.OrderItem

		for _, item := range orderItemPack.OrderItems {
			var food models.Food
			err := foodCollection.FindOne(ctx, bson.M{"foodId": item.FoodID}).Decode(&food)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Food item not found", "foodId": item.FoodID})
				return
			}

			unitPrice := food.Price

			orderItem := models.OrderItem{
				ID:          primitive.NewObjectID(),
				Quantity:    &item.Quantity,
				UnitPrice:   unitPrice,
				CreatedAt:   time.Now().UTC(),
				UpdatedAt:   time.Now().UTC(),
				FoodID:      &item.FoodID,
				OrderItemID: primitive.NewObjectID().Hex(),
				OrderID:     orderId,
			}

			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)
			createdOrderItems = append(createdOrderItems, orderItem)
		}

		_, err = orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order items"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order and items created successfully", "totalCount": len(createdOrderItems), "orderItems": createdOrderItems})
	}
}

func GetAllOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := orderItemCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while retrieving order items"})
			return
		}
		defer cursor.Close(ctx)

		var allOrderItems []models.OrderItem
		if err := cursor.All(ctx, &allOrderItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while decoding order items"})
			return
		}

		if len(allOrderItems) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"message": "No order items found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"totalCount": len(allOrderItems), "orderItems": allOrderItems})
	}
}

func GetOrderItemByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderItemId := c.Param("orderItemId")
		var orderItem models.OrderItem

		err := orderItemCollection.FindOne(ctx, bson.M{"orderItemId": orderItemId}).Decode(&orderItem)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurred while retrieving the order item"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"orderItem": orderItem})
	}
}

func GetOrderItemsByOrderID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderId := c.Param("orderId")

		cursor, err := orderItemCollection.Find(ctx, bson.M{"orderId": orderId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order items"})
			return
		}
		defer cursor.Close(ctx)

		var orderItems []models.OrderItem
		if err := cursor.All(ctx, &orderItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while parsing order items for order"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"totalCount": len(orderItems), "orderItems": orderItems})
	}
}

func UpdateOrderItemByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItem models.OrderItem
		orderItemId := c.Param("orderItemId")

		if err := c.BindJSON(&orderItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		filter := bson.M{"orderItemId": orderItemId}
		updateObj := bson.D{}

		if orderItem.UnitPrice != nil {
			updateObj = append(updateObj, bson.E{Key: "unitprice", Value: *orderItem.UnitPrice})
		}

		if orderItem.Quantity != nil {
			updateObj = append(updateObj, bson.E{Key: "quantity", Value: *orderItem.Quantity})
		}

		if orderItem.FoodID != nil {
			updateObj = append(updateObj, bson.E{Key: "foodId", Value: *orderItem.FoodID})
		}

		orderItem.UpdatedAt = time.Now().UTC()
		updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: orderItem.UpdatedAt})

		if len(updateObj) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
			return
		}

		opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)

		var updatedOrderItem models.OrderItem
		err := orderItemCollection.FindOneAndUpdate(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, opts).Decode(&updatedOrderItem)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order item"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order item updated successfully", "orderItem": updatedOrderItem})
	}
}
