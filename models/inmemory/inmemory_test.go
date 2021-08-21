package inmemory_test

import (
	"fmt"
	"testing"

	"github.com/ocakhasan/getir-api-task/controllers/responses"

	"github.com/ocakhasan/getir-api-task/models/inmemory"
)

func TestInMemory_Get(t *testing.T) {
	key := "getir"
	value := "company"
	mockMemory := inmemory.New(map[string]string{
		key: value,
	})

	// test existing key
	got, ok := mockMemory.Get(key)
	if got != value || !ok {
		t.Errorf("expected: %s, got: %s\n", value, got)
	}

	// test not existing key
	got, ok = mockMemory.Get("non-existing-key")
	fmt.Println(got, ok)
	if got != "" || ok {
		t.Errorf("expected: , got: %v\n", got)
	}
}

func TestInMemory_Set(t *testing.T) {
	mockMemory := inmemory.New(map[string]string{})
	key := "getir"
	value := "company"
	mockMemory.Set(*responses.NewInMemoryBody(key, value))

	got, ok := mockMemory.Get(key)
	if got != value || !ok {
		t.Errorf("Expected %s: got: %s\n", value, got)
	}

}
