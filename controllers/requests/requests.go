// Package requests handles all structures related with HTTP responses
package requests

// MongoRequestBody represents request body structure for MongoDB POST operations.
type MongoRequestBody struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}
