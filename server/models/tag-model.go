package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tag model
type Tag struct {
	ID     primitive.ObjectID `bson:"_id"`
	PostId *string            `bson:"post_id" json:"post_id" validate:"required"`
	Name   *string            `bson:"name" json:"name" validate:"required"`
}
