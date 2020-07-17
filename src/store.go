package src

import "fmt"

// CollectorStoreIteratorFunc callback
type CollectorStoreIteratorFunc func(string, []byte)

// CollectorStore interface
type CollectorStore interface {
	Put(string, []byte)
	Get(string) ([]byte, error)
	Delete(string) error
	IterateAll(CollectorStoreIteratorFunc)
}

// SimpleInMemStore example implementation. For learning purpose only. Its not scaleable
type SimpleInMemStore struct {
	CollectorStore
	cache map[string][]byte
}

// CreateSimpleInMemStore factory
func CreateSimpleInMemStore() CollectorStore {
	return &SimpleInMemStore{
		cache: map[string][]byte{},
	}
}

// Put implementation
func (sims *SimpleInMemStore) Put(key string, data []byte) {
	sims.cache[key] = data
	fmt.Printf("Store put %s\n%s\n", key, string(data))
}

// Get implementation
func (sims *SimpleInMemStore) Get(key string) ([]byte, error) {
	if data, ok := sims.cache[key]; ok {
		return data, nil
	}
	return nil, fmt.Errorf("Entity not found '%s'", key)
}

// Delete implementation
func (sims *SimpleInMemStore) Delete(key string) error {
	if _, ok := sims.cache[key]; ok {
		delete(sims.cache, key)
	}

	return nil
}

// IterateAll implementation
func (sims *SimpleInMemStore) IterateAll(iterator CollectorStoreIteratorFunc) {
	for k, v := range sims.cache {
		iterator(k, v)
	}
}
