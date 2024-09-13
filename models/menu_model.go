package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Name      string             `json:"name" validate:"required"`
	Category  string             `json:"category" validate:"required"`
	StartDate *time.Time         `json:"startDate" bson:"startDate"`
	EndDate   *time.Time         `json:"endDate" bson:"endDate"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	MenuID    string             `json:"menuId" bson:"menuId"`
}
