// Package controllers handles all of HTTP requests and responses
package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ocakhasan/getir-api-task/controllers/errors"

	"github.com/ocakhasan/getir-api-task/controllers/requests"

	modelMongo "github.com/ocakhasan/getir-api-task/models/mongo"

	"github.com/ocakhasan/getir-api-task/controllers/responses"

	"github.com/ocakhasan/getir-api-task/models/inmemory"
)

// contentType is application/json for our program.
const contentType = "application/json"

// RequestHandler is a structure which handles all of request operations
type RequestHandler struct {
	inMemoryDB inmemory.InMemory
	mongoDB    *modelMongo.DB
}

// New returns a new RequestHandler with given mongoDB and inMemory structures.
func New(mongoDB *modelMongo.DB, inmemoryDB inmemory.InMemory) *RequestHandler {
	return &RequestHandler{
		inMemoryDB: inmemoryDB,
		mongoDB:    mongoDB,
	}
}

// GetMongo handles POST operations for mongoDB related requests.
func (h *RequestHandler) GetMongo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", contentType)
	if r.Method != http.MethodPost {
		errors.WriteError(w, errors.ErrorUnsupportedMethod)
		return
	}
	if r.Header.Get("content-type") != contentType {
		fmt.Println("I am here content-type")
		errors.WriteError(w, errors.ErrorContentType)
		return
	}

	// decode body to filter.
	filter := requests.MongoRequestBody{}
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		errors.WriteError(w, errors.ErrorInvalidBody)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			errors.WriteError(w, err)
			return
		}
	}(r.Body)

	// create pipe from filter (request body)
	pipe, err := modelMongo.CreatePipe(filter)
	if err != nil {
		fmt.Println("pipe error", err)
		errors.WriteError(w, err)
		return
	}

	// get the results with pipe.
	ctx := context.Background()
	results, err := h.mongoDB.Get(ctx, pipe)
	if err != nil {
		fmt.Println("aggregte", err)
		errors.WriteError(w, err)
		return
	}

	response := responses.NewMongoResponse(0, "success", results)
	if err := json.NewEncoder(w).Encode(&response); err != nil {
		errors.WriteError(w, err)
		return
	}
}

// InMemory handles GET and POST related requests about in-memory operations.
func (h *RequestHandler) InMemory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", contentType)
	if r.Method == http.MethodPost { // Post Method
		if r.Header.Get("content-type") != contentType {
			errors.WriteError(w, errors.ErrorContentType)
			return
		}
		requestBody := &responses.InMemoryBody{}
		if err := json.NewDecoder(r.Body).Decode(requestBody); err != nil {
			errors.WriteError(w, errors.ErrorInvalidBody)
			return
		}
		h.inMemoryDB.Set(*requestBody)
		w.WriteHeader(http.StatusCreated)
	} else if r.Method == http.MethodGet { // Get Method
		key := r.URL.Query().Get("key")
		value, ok := h.inMemoryDB.Get(key)
		if !ok {
			errors.WriteError(w, errors.ErrorNotFound)
			return
		}

		responseBody := responses.NewInMemoryBody(key, value)
		if err := json.NewEncoder(w).Encode(responseBody); err != nil {
			errors.WriteError(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else { // Neither Post nor Get
		errors.WriteError(w, errors.ErrorUnsupportedMethod)
		return
	}
}
