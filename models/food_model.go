package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      *string            `json:"name" validate:"required,min=2,max=100"`
	Price     *float64           `json:"price" validate:"required"`
	FoodImage *string            `json:"foodImage" bson:"foodImage" validate:"required"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	FoodID    string             `json:"foodId" bson:"foodId"`
	MenuID    *string            `json:"menuId" bson:"menuId" validate:"required"`
}
