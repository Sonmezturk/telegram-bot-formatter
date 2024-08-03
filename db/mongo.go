package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
  
var (
	clientInstance *mongo.Client
	clientInstanceError error
	mongoOnce sync.Once
)

func MongoInit() {
	mongoConnect()
}

func mongoConnect() {
	mongoStr := os.Getenv("MONGO_STR")
	mongoOnce.Do(func() {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(mongoStr).SetServerAPIOptions(serverAPI)

	var err error
	clientInstance, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		clientInstanceError = err
		log.Fatal(err)
	}

	err = clientInstance.Ping(context.TODO(), nil)
	if err != nil {
		clientInstanceError = err
		log.Fatal(err)
	}

	fmt.Println("Successfully connected and pinged MongoDB.")
})
}

func GetMongoClient() (*mongo.Client, error) {
	return clientInstance, clientInstanceError
}