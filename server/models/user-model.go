package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model
type User struct {
	ID               primitive.ObjectID `bson:"_id"`
	FirstName        *string            `bson:"first_name" json:"first_name" validate:"required,min=2,max=100"`
	LastName         *string            `bson:"last_name" json:"last_name" validate:"required,min=2,max=100"`
	Email            *string            `bson:"email" json:"email" validate:"email,required"`
	Password         *string            `bson:"password" json:"password" validate:"required,min=6"`
	UserType         *string            `bson:"user_type"`
	IsEmailVerified  bool               `bson:"is_email_verified" json:"is_email_verified"`
	VerificationCode *string            `bson:"verification_code" json:"verification_code"`
	Token            *string            `bson:"token" json:"token"`
	RefreshToken     *string            `bson:"refresh_token" json:"refresh_token"`
	CreatedAt        time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at" json:"updated_at"`
}
