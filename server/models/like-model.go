package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Like model
type Like struct {
	ID     primitive.ObjectID `bson:"id"`
	UserId *string            `bson:"user_id" json:"user_id"`
	PostId *string            `bson:"post_id" json:"post_id"`
}
