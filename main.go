package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/ocakhasan/getir-api-task/models/inmemory"

	mongo2 "go.mongodb.org/mongo-driver/mongo"

	"github.com/ocakhasan/getir-api-task/models/mongo"

	"github.com/ocakhasan/getir-api-task/controllers"
)

func main() {

	mongoClient, collection := mongo.NewClient()
	defer func(mongoClient *mongo2.Client, ctx context.Context) {
		err := mongoClient.Disconnect(ctx)
		if err != nil {
			log.Fatalf("Disconnect error %v\n", err)
		}
	}(mongoClient, context.Background())

	mongoModel := mongo.New(collection)
	inMemory := inmemory.New(map[string]string{})

	agent := controllers.New(mongoModel, inMemory)
	mux := http.NewServeMux()
	mux.HandleFunc("/inmemory", agent.InMemory)
	mux.HandleFunc("/records", agent.GetMongo)

	address := fmt.Sprintf(":%d", 3000)
	log.Fatal(http.ListenAndServe(address, mux))
}
