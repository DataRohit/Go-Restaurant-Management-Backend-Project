package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	OrderDate time.Time          `json:"orderDate" bson:"orderDate" validate:"required"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	OrderID   string             `json:"orderId" bson:"orderId"`
	TableID   *string            `json:"tableId" bson:"tableId" validate:"required"`
}
