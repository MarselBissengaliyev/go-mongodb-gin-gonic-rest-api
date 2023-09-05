package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Comment model
type Comment struct {
	ID      primitive.ObjectID `bson:"id"`
	PostId  *string            `bson:"post_id" json:"post_id"`
	UserId  *string            `bson:"user_id" json:"user_id"`
	Content *string            `bson:"content" json:"content"`
}
