package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/ocakhasan/getir-api-task/controllers"
	"github.com/ocakhasan/getir-api-task/models/inmemory"
	modelMongo "github.com/ocakhasan/getir-api-task/models/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mongoClient, collection := modelMongo.NewClient()
	defer func(mongoClient *mongo.Client, ctx context.Context) {
		err := mongoClient.Disconnect(ctx)
		if err != nil {
			log.Fatalf("Disconnect error %v\n", err)
		}
	}(mongoClient, context.Background())

	mongoModel := modelMongo.New(collection)
	inMemory := inmemory.New(map[string]string{})

	agent := controllers.New(mongoModel, inMemory)
	mux := http.NewServeMux()
	mux.HandleFunc("/inmemory", agent.InMemory)
	mux.HandleFunc("/records", agent.GetMongo)

	PORT := os.Getenv("PORT")
	address := fmt.Sprintf(":%s", PORT)
	log.Fatal(http.ListenAndServe(address, mux))
}
