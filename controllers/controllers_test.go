package controllers_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ocakhasan/getir-api-task/controllers/requests"

	"github.com/ocakhasan/getir-api-task/controllers/responses"

	"github.com/ocakhasan/getir-api-task/controllers"
	"github.com/ocakhasan/getir-api-task/models/inmemory"
	mongoModel "github.com/ocakhasan/getir-api-task/models/mongo"
)

func newTestServer() (*httptest.Server, mongo.Client) {

	mongoClient, collection := mongoModel.NewClient()

	// create mongo db
	mongoDB := mongoModel.New(collection)

	// create in memory
	inMemory := inmemory.New(map[string]string{
		"getir": "company",
	})

	agent := controllers.New(mongoDB, inMemory)
	mux := http.NewServeMux()
	mux.HandleFunc("/inmemory", agent.InMemory)
	mux.HandleFunc("/records", agent.GetMongo)

	return httptest.NewServer(mux), *mongoClient
}

func TestRequestHandler_GetMongo(t *testing.T) {

	ts, mongoClient := newTestServer()
	defer func(mongoClient *mongo.Client, ctx context.Context) {
		err := mongoClient.Disconnect(ctx)
		if err != nil {
			log.Fatalf("Disconnect error %v\n", err)
		}
	}(&mongoClient, context.Background())

	defer ts.Close()
	client := ts.Client()

	cases := []struct {
		name        string
		httpMethod  string
		requestBody requests.MongoRequestBody
		statusCode  int
		contentType string
	}{
		{
			name:       "normal post",
			httpMethod: http.MethodPost,
			requestBody: requests.MongoRequestBody{
				StartDate: "2020-01-02",
				EndDate:   "2021-01-02",
				MinCount:  2700,
				MaxCount:  3000,
			},
			statusCode:  http.StatusOK,
			contentType: "application/json",
		},
		{
			name:       "unsupported Method",
			httpMethod: http.MethodGet,
			requestBody: requests.MongoRequestBody{
				StartDate: "2020-01-02",
				EndDate:   "2021-01-02",
				MinCount:  2700,
				MaxCount:  3000,
			},
			statusCode:  http.StatusMethodNotAllowed,
			contentType: "application/json",
		},
		{
			name:       "false invalid body enddate smaller",
			httpMethod: http.MethodPost,
			requestBody: requests.MongoRequestBody{
				StartDate: "2020-01-02",
				EndDate:   "2019-01-02",
				MinCount:  2700,
				MaxCount:  3000,
			},
			statusCode:  http.StatusBadRequest,
			contentType: "application/json",
		},
		{
			name:       "unsupported application type",
			httpMethod: http.MethodPost,
			requestBody: requests.MongoRequestBody{
				StartDate: "2020-01-02",
				EndDate:   "2019-01-02",
				MinCount:  2700,
				MaxCount:  3000,
			},
			statusCode:  http.StatusBadRequest,
			contentType: "application/xml",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/records", ts.URL)
			fmt.Println("url", url)

			jsonBody, _ := json.Marshal(tt.requestBody)
			request, err := http.NewRequest(tt.httpMethod, url, strings.NewReader(string(jsonBody)))
			if err != nil {
				t.Errorf("unexpected error %v\n", err)
			}

			request.Header.Set("content-type", tt.contentType)
			request.Header.Set("Referer", "http://localhost")

			response, err := client.Do(request)
			if err != nil {
				t.Errorf("unexpected response error %v\n", err)
			}

			if response.StatusCode != tt.statusCode {
				t.Errorf("expected statusCode :%v, got: %v\n", tt.statusCode, response.StatusCode)
			}
		})
	}

}

func TestRequestHandler_InMemory(t *testing.T) {
	ts, _ := newTestServer()
	defer ts.Close()
	client := ts.Client()

	cases := []struct {
		name        string
		httpMethod  string
		requestBody responses.InMemoryBody
		statusCode  int
		contentType string
		key         string
	}{
		{
			name:        "get a valid data (already stored)",
			httpMethod:  http.MethodGet,
			statusCode:  http.StatusOK,
			contentType: "application/json",
			key:         "getir",
		},
		{
			name:        "get an invalid data",
			httpMethod:  http.MethodGet,
			statusCode:  http.StatusNotFound,
			contentType: "application/json",
			key:         "invalidKey",
		},
		{
			name:        "post a key value",
			httpMethod:  http.MethodPost,
			requestBody: *responses.NewInMemoryBody("hasan", "ocak"),
			statusCode:  http.StatusCreated,
			contentType: "application/json",
		},
		{
			name:        "unsupported http method",
			httpMethod:  http.MethodDelete,
			statusCode:  http.StatusMethodNotAllowed,
			contentType: "application/json",
		},
		{
			name:        "get a valid data (new added)",
			httpMethod:  http.MethodGet,
			statusCode:  http.StatusOK,
			contentType: "application/json",
			key:         "hasan",
		},
		{
			name:        "false content type",
			httpMethod:  http.MethodPost,
			statusCode:  http.StatusBadRequest,
			contentType: "application/xml",
			key:         "hasan",
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			url := fmt.Sprintf("%s/inmemory?key=%s", ts.URL, tt.key)

			jsonBody, _ := json.Marshal(tt.requestBody)
			request, err := http.NewRequest(tt.httpMethod, url, strings.NewReader(string(jsonBody)))
			if err != nil {
				t.Errorf("unexpected error %v\n", err)
			}

			request.Header.Set("content-type", tt.contentType)
			request.Header.Set("Referer", "http://localhost")

			response, err := client.Do(request)
			if err != nil {
				t.Errorf("unexpected response error %v\n", err)
			}

			if response.StatusCode != tt.statusCode {
				t.Errorf("expected statusCode :%v, got: %v\n", tt.statusCode, response.StatusCode)
			}
		})
	}
}
