package helpers

import (
	"context"
	"log"
	"time"

	"github.com/MarselBisengaliev/go-react-blog/config"
	"github.com/MarselBisengaliev/go-react-blog/database"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TokenHelper struct{}

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	UserType  string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

// generate all tokens helper
func (h *TokenHelper) GenerateAllTokens(
	email string,
	firstName string,
	lastName string,
	userType string,
	uid string) (signedToken string, signedRefreshToken string, err error) {

	conf, err := config.LoadConfig("./")

	if err != nil {
		log.Fatal("could not load config", err)
		return
	}

	claims := &SignedDetails{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		UserType:  userType,
		Uid:       uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(conf.TokenSecret))
	if err != nil {
		log.Panic(err)
		return
	}

	refreshToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		refreshClaims,
	).SignedString([]byte(conf.TokenSecret))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

// validate token helper
func (h *TokenHelper) ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	conf, err := config.LoadConfig("./")

	if err != nil {
		log.Fatal("could not load config", err)
		return
	}

	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(conf.TokenSecret), nil
		},
	)

	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}

	return claims, msg
}

// update all tokens helper
func (h *TokenHelper) UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	UpdatedAt, parseErr := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	if parseErr != nil {
		log.Panic(parseErr)
		return
	}

	var updateObj = bson.M{
		"token":         signedToken,
		"refresh_token": signedRefreshToken,
		"updated_at":    UpdatedAt,
	}

	userObjectId, err := primitive.ObjectIDFromHex(userId)

	if err != nil {
		log.Panic(err)
		return
	}

	filter := bson.M{"_id": userObjectId}
	upsert := true
	opts := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, updateOneErr := userCollection.UpdateOne(ctx, filter, bson.M{"$set": updateObj}, &opts)
	if updateOneErr != nil {
		log.Panic(err)
		return
	}
}
