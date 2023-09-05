package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenModel struct {
	ID           primitive.ObjectID `bson:"_id"`
	Token        *string            `bson:"token" json:"token"`
	RefreshToken *string            `bson:"refresh_token" json:"refresh_token"`
	UserId       *string            `bson:"user_id" json:"user_id"`
	UserAgent    *string            `bson:"device" json:"device"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}
