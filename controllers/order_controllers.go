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
	"go.mongodb.org/mongo-driver/mongo/options"
)

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "order")

func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var order models.Order
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}

		if err := validate.Struct(order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if *(order.TableID) != "" {
			var table models.Table
			err := tableCollection.FindOne(ctx, bson.M{"tableId": order.TableID}).Decode(&table)
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
				return
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify table"})
				return
			}
		}

		order.CreatedAt = time.Now().UTC()
		order.UpdatedAt = time.Now().UTC()
		order.ID = primitive.NewObjectID()
		order.OrderID = order.ID.Hex()

		_, err := orderCollection.InsertOne(ctx, order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "order": order})
	}
}

func GetAllOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := orderCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
			return
		}
		defer cursor.Close(ctx)

		var orders []models.Order
		if err := cursor.All(ctx, &orders); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while parsing orders"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"totalCount": len(orders), "orders": orders})
	}
}

func GetOrderByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		orderID := c.Param("orderId")
		if orderID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
			return
		}

		var order models.Order
		err := orderCollection.FindOne(ctx, bson.M{"orderId": orderID}).Decode(&order)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"order": order})
	}
}

func UpdateOrderByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		orderID := c.Param("orderId")
		if orderID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order ID is required"})
			return
		}

		var order models.Order
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}

		updateObj := bson.D{}
		if *(order.TableID) != "" {
			var table models.Table
			err := tableCollection.FindOne(ctx, bson.M{"tableId": order.TableID}).Decode(&table)
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
				return
			} else if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify table"})
				return
			}
			updateObj = append(updateObj, bson.E{Key: "tableId", Value: order.TableID})
		}

		order.UpdatedAt = time.Now().UTC()
		updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: order.UpdatedAt})

		if len(updateObj) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
			return
		}

		filter := bson.M{"orderId": orderID}
		opts := options.Update().SetUpsert(false)

		_, err := orderCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
			return
		}

		var updatedOrder models.Order
		err = orderCollection.FindOne(ctx, filter).Decode(&updatedOrder)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated order"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order updated successfully", "order": updatedOrder})
	}
}
