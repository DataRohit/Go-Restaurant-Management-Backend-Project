package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Table struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	NumberOfGuests *int               `json:"numberOfGuests" bson:"numberOfGuests" validate:"required"`
	TableNumber    *int               `json:"tableNumber" bson:"tableNumber" validate:"required"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt" bson:"updatedAt"`
	TableID        string             `json:"tableId" bson:"tableId"`
}
