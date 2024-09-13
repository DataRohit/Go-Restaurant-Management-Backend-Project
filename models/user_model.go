package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	FirstName    *string            `json:"firstName" bson:"firstName" validate:"required,min=2,max=100"`
	LastName     *string            `json:"lastName" bson:"lastName" validate:"required,min=2,max=100"`
	Password     *string            `json:"password" validate:"required,min=6"`
	Email        *string            `json:"email" validate:"email,required"`
	Avatar       *string            `json:"avatar"`
	Phone        *string            `json:"phone" validate:"required"`
	AccessToken  *string            `json:"accessToken" bson:"accessToken"`
	RefreshToken *string            `json:"refreshToken" bson:"refreshToken"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
	UserID       string             `json:"userId" bson:"userId"`
}
