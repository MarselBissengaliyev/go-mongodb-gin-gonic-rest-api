package controllers

import (
	"github.com/MarselBisengaliev/go-react-blog/database"
	"github.com/MarselBisengaliev/go-react-blog/helpers"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var postCollection *mongo.Collection = database.OpenCollection(database.Client, "posts")
var commentCollection *mongo.Collection = database.OpenCollection(database.Client, "comments")
var tagCollection *mongo.Collection = database.OpenCollection(database.Client, "tags")
var likeCollection *mongo.Collection = database.OpenCollection(database.Client, "likes")
var authHelper = new(helpers.AuthHelper)
var tokenHelper = new(helpers.TokenHelper)
var mongoHelper = new(helpers.MongoHelper)
var tagHelper = new(helpers.TagHelper)
var encodeHelper = new(helpers.EncodeHelper)
var emailHelper = new(helpers.EmailHelper)
