package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	InvoiceID      string             `json:"invoiceId" bson:"invoiceId"`
	OrderID        string             `json:"orderId" bson:"orderId" validate:"required"`
	PaymentMethod  *string            `json:"paymentMethod" bson:"paymentMethod" validate:"eq=CARD|eq=CASH|eq=ONLINE"`
	PaymentStatus  *string            `json:"paymentStatus" bson:"paymentStatus" validate:"required,eq=PENDING|eq=PAID"`
	PaymentDueDate time.Time          `json:"paymentDueDate" bson:"paymentDueDate" validate:"required"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt" bson:"updatedAt"`
}
