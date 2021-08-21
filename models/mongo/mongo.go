// Package mongo handles all of database operation with MongoDB.
package mongo

import (
	"context"
	"log"
	"time"

	"github.com/ocakhasan/getir-api-task/controllers/errors"

	"github.com/ocakhasan/getir-api-task/controllers/requests"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo/options"

	driver "go.mongodb.org/mongo-driver/mongo"
)

// Config is to create everything related with MONGO
type Config struct {
	URI        string
	Database   string
	Collection string
}

// Object represents MongoDB Objects after filter operation.
type Object struct {
	Key        string    `json:"key" bson:"key"`
	Val        string    `json:"val" bson:"val"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	Counts     []int64   `json:"counts" bson:"counts"`
	TotalCount int64     `json:"totalCount" bson:"totalCount"`
}

// DB represents MongoDB interface for the request handling operations
type DB struct {
	Collection *driver.Collection
}

// New returns a new DB with given collection.
func New(collection *driver.Collection) *DB {
	return &DB{
		Collection: collection,
	}
}

// NewClient creates mongo client and returns given collection.
func NewClient() (*driver.Client, *driver.Collection) {
	config := CreateConfig()
	clientOptions := options.Client().ApplyURI(config.URI)
	mongoClient, err := driver.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Error while creating mongo client :%v\n", err)
	}
	ctx := context.Background()
	err = mongoClient.Connect(ctx)
	if err != nil {
		log.Fatalf("Error while connecting to mongo client :%v\n", err)
	}
	collection := mongoClient.Database(config.Database).Collection(config.Collection)
	return mongoClient, collection
}

// CreatePipe creates a pipe structure with given request body.
func CreatePipe(filter requests.MongoRequestBody) ([]bson.M, error) {
	if filter.MaxCount < filter.MinCount {
		return nil, errors.ErrorInvalidBody
	}
	startDate, err := getTime(filter.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := getTime(filter.EndDate)
	if err != nil {
		return nil, err
	}

	if startDate.After(endDate) {
		return nil, errors.ErrorInvalidBody
	}

	pipeLine := []bson.M{{"$project": bson.M{
		"key":        1,
		"createdAt":  1,
		"totalCount": bson.M{"$sum": "$counts"},
	}}, {"$match": bson.M{
		"$and": []bson.M{
			{"createdAt": bson.M{"$gt": startDate}},
			{"createdAt": bson.M{"$lt": endDate}},
			{"totalCount": bson.M{"$gt": filter.MinCount}},
			{"totalCount": bson.M{"$lt": filter.MaxCount}}},
	}},
	}
	return pipeLine, nil
}

// Get returns all of the objects which fits in pipe.
func (db *DB) Get(ctx context.Context, pipe []bson.M) ([]Object, error) {
	cursor, err := db.Collection.Aggregate(ctx, pipe)
	if err != nil {
		return nil, err
	}
	defer func(cursor *driver.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(cursor, ctx)

	var results []Object
	for cursor.Next(ctx) {
		o := Object{}
		if err := cursor.Decode(&o); err != nil {
			return nil, err
		}
		results = append(results, o)
	}
	return results, nil
}

// getTime is used to parse filter time, it checks if given time is in right format
func getTime(s string) (time.Time, error) {
	parsedTime, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}

// CreateConfig returns config values related with MongoDB
func CreateConfig() Config {
	return Config{
		URI:        "mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getir\u0002case-study?retryWrites=true",
		Database:   "getir-case-study",
		Collection: "records",
	}
}
