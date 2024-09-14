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

type InvoiceViewFormat struct {
	InvoiceID      string
	PaymentMethod  string
	OrderID        string
	PaymentStatus  *string
	TableNumber    interface{}
	PaymentDueDate time.Time
	OrderDetails   interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var invoice models.Invoice

		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		var order models.Order
		err := orderCollection.FindOne(ctx, bson.M{"orderId": invoice.OrderID}).Decode(&order)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		status := "PENDING"
		if invoice.PaymentStatus == nil {
			invoice.PaymentStatus = &status
		}

		now := time.Now().UTC()
		invoice.PaymentDueDate = now.AddDate(0, 0, 1)
		invoice.CreatedAt = now
		invoice.UpdatedAt = now
		invoice.ID = primitive.NewObjectID()
		invoice.InvoiceID = invoice.ID.Hex()

		if err := validate.Struct(invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed: " + err.Error()})
			return
		}

		_, err = invoiceCollection.InsertOne(ctx, invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invoice"})
			return
		}

		var createdInvoice models.Invoice
		err = invoiceCollection.FindOne(ctx, bson.M{"invoiceId": invoice.InvoiceID}).Decode(&createdInvoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created invoice"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Invoice created successfully", "invoice": createdInvoice})
	}
}

func GetAllInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := invoiceCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invoices"})
			return
		}

		var allInvoices []models.Invoice
		if err := cursor.All(ctx, &allInvoices); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse invoices"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"invoices": allInvoices, "totalCount": len(allInvoices)})
	}
}

func GetInvoiceByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		invoiceID := c.Param("invoiceId")

		var invoice models.Invoice
		err := invoiceCollection.FindOne(ctx, bson.M{"invoiceId": invoiceID}).Decode(&invoice)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve invoice"})
			}
			return
		}

		var order models.Order
		err = orderCollection.FindOne(ctx, bson.M{"orderId": invoice.OrderID}).Decode(&order)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve order"})
			return
		}

		var table models.Table
		err = tableCollection.FindOne(ctx, bson.M{"tableId": order.TableID}).Decode(&table)
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve table"})
			return
		}

		cursor, err := orderItemCollection.Find(ctx, bson.M{"orderId": invoice.OrderID})
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

		if len(orderItems) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No order items found for this invoice"})
			return
		}

		invoiceView := InvoiceViewFormat{
			OrderID:        invoice.OrderID,
			PaymentDueDate: invoice.PaymentDueDate,
			PaymentMethod:  helper.GetNonNilString(invoice.PaymentMethod, "null"),
			InvoiceID:      invoice.InvoiceID,
			PaymentStatus:  invoice.PaymentStatus,
			TableNumber:    table.TableNumber,
			OrderDetails:   orderItems,
		}

		c.JSON(http.StatusOK, gin.H{"invoice": invoiceView})
	}
}

func UpdateInvoiceByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		invoiceID := c.Param("invoiceId")
		var updateData models.Invoice

		if err := c.BindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		filter := bson.M{"invoiceId": invoiceID}
		updateFields := bson.D{}

		if updateData.PaymentMethod != nil {
			updateFields = append(updateFields, bson.E{Key: "paymentMethod", Value: updateData.PaymentMethod})
		}

		if updateData.PaymentStatus != nil {
			updateFields = append(updateFields, bson.E{Key: "paymentStatus", Value: updateData.PaymentStatus})
		} else {
			defaultStatus := "PENDING"
			updateFields = append(updateFields, bson.E{Key: "paymentStatus", Value: defaultStatus})
		}

		updateFields = append(updateFields, bson.E{Key: "updatedAt", Value: time.Now().UTC()})

		update := bson.D{{Key: "$set", Value: updateFields}}

		result, err := invoiceCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invoice: " + err.Error()})
			return
		}

		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
			return
		}

		var updatedInvoice models.Invoice
		err = invoiceCollection.FindOne(ctx, filter).Decode(&updatedInvoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated invoice: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Invoice updated successfully", "invoice": updatedInvoice})
	}
}
