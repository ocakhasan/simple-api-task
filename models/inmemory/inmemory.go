// Package inmemory handles all of database operations for in memory
package inmemory

import (
	"sync"

	"github.com/ocakhasan/getir-api-task/controllers/responses"
)

// InMemory represents in memory data structure for request handling.
type InMemory struct {
	db map[string]string
	mu *sync.RWMutex
}

// New returns a new InMemory object
func New(db map[string]string) InMemory {
	return InMemory{
		db: db,
		mu: &sync.RWMutex{},
	}
}

// Get returns corresponding string value for given key
func (m *InMemory) Get(key string) (string, bool) {
	m.mu.RLock() // use R lock for reading, so others can also read.
	defer m.mu.RUnlock()
	val, ok := m.db[key]
	if !ok {
		return "", false
	}
	return val, ok
}

// Set writes the given key-value pair to database.
func (m *InMemory) Set(body responses.InMemoryBody) {
	m.mu.Lock() // use normal lock
	defer m.mu.Unlock()
	m.db[body.Key] = body.Value
}
