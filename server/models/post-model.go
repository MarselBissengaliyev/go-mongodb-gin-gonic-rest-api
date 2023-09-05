package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post model
type Post struct {
	ID        primitive.ObjectID `bson:"id"`
	UserId    *string            `json:"user_id" validate:"required"`
	Content   *string            `bson:"content" json:"content" validate:"required"`
	Title     *string            `bson:"title" json:"title" validate:"required"`
	Preview   *string            `bson:"preview" json:"preview" validate:"required"`
	Views     *int               `bson:"views" json:"views" validate:"required"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	Tags      []Tag              `json:"tags" validate:"required"`
}
