package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderItem struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Quantity    *string            `json:"quantity" bson:"quantity" validate:"required,eq=S|eq=M|eq=L"`
	UnitPrice   *float64           `json:"unitPrice" bson:"unitPrice" validate:"required"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
	FoodId      *string            `json:"foodId" bson:"foodId" validate:"required"`
	OrderItemId string             `json:"orderItemId" bson:"orderItemId"`
	OrderId     string             `json:"orderId" bson:"orderId" validate:"required"`
}
