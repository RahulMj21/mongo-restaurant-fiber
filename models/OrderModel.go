package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	OrderDate time.Time          `json:"order_date" validate:"required"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	TableId   *string            `json:"table_id"`
}
