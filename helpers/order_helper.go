package helpers

import (
	"context"
	"time"

	"github.com/datarohit/go-restaurant-management-backend-project/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func OrderItemOrderCreator(order models.Order, ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	order.CreatedAt = time.Now().UTC()
	order.UpdatedAt = time.Now().UTC()
	order.ID = primitive.NewObjectID()
	order.OrderID = order.ID.Hex()

	_, err := orderCollection.InsertOne(ctx, order)
	if err != nil {
		return "", err
	}

	return order.OrderID, nil
}
