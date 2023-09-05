package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MarselBisengaliev/go-react-blog/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func DBinstance() *mongo.Client {
	conf, err := config.LoadConfig("./")

	if err != nil {
		log.Fatal("could not load config", err)
	}

	MongoUri := conf.DBHost
	fmt.Print(MongoUri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoUri))

	if err != nil {
		log.Fatal(err)
	}

	defer cancel()
	fmt.Println("Connected to MongoDB")
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("blogdb").Collection(collectionName)
	return collection
}
