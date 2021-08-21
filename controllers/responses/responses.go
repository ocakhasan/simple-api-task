// Package responses handles all structures related with HTTP responses
package responses

import "github.com/ocakhasan/getir-api-task/models/mongo"

// InMemoryBody represents body structure for in-memory POST operation
type InMemoryBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewInMemoryBody returns a new InMemoryBody with given values.
func NewInMemoryBody(key, value string) *InMemoryBody {
	return &InMemoryBody{
		Key:   key,
		Value: value,
	}
}

// MongoResponse represents response body for the mongoDB related requests.
type MongoResponse struct {
	Code    int            `json:"code"`
	Msg     string         `json:"msg"`
	Records []mongo.Object `json:"records"`
}

// NewMongoResponse returns a new MongoResponse with given values.
func NewMongoResponse(code int, msg string, records []mongo.Object) *MongoResponse {
	return &MongoResponse{
		Code:    code,
		Msg:     msg,
		Records: records,
	}
}
