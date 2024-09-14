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

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func CreateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var table models.Table
		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}

		if validationErr := validate.Struct(table); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		table.CreatedAt = time.Now().UTC()
		table.UpdatedAt = time.Now().UTC()
		table.ID = primitive.NewObjectID()
		table.TableID = table.ID.Hex()

		_, err := tableCollection.InsertOne(ctx, table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create table"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Table created successfully", "table": table})
	}
}

func GetAllTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := tableCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tables"})
			return
		}
		defer cursor.Close(ctx)

		var tables []models.Table
		if err := cursor.All(ctx, &tables); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while parsing tables"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"totalCount": len(tables), "tables": tables})
	}
}

func GetTableByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tableID := c.Param("tableId")
		if tableID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Table ID is required"})
			return
		}

		var table models.Table
		err := tableCollection.FindOne(ctx, bson.M{"tableId": tableID}).Decode(&table)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve table"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"table": table})
	}
}

func UpdateTableByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		tableID := c.Param("tableId")
		if tableID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Table ID is required"})
			return
		}

		var table models.Table
		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}

		updateObj := bson.D{}
		if table.NumberOfGuests != nil {
			updateObj = append(updateObj, bson.E{Key: "numberOfGuests", Value: table.NumberOfGuests})
		}

		if table.TableNumber != nil {
			updateObj = append(updateObj, bson.E{Key: "tableNumber", Value: table.TableNumber})
		}

		table.UpdatedAt = time.Now().UTC()
		updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: table.UpdatedAt})

		if len(updateObj) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
			return
		}

		filter := bson.M{"tableId": tableID}
		opts := options.Update().SetUpsert(false)

		result, err := tableCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, opts)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update table"})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
			return
		}

		var updatedTable models.Table
		err = tableCollection.FindOne(ctx, filter).Decode(&updatedTable)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated table"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Table updated successfully", "table": updatedTable})
	}
}
