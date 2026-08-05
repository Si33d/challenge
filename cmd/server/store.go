package main

import "sync"

type Store struct {
	data map[string]string
}

var mu sync.RWMutex

var store = NewStore() //level3 done here only

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

func (s *Store) Set(key, value string) {
	mu.Lock()
	defer mu.Unlock()
	s.data[key] = value

}

func (s *Store) Get(key string) (string, bool) {
	mu.RLock()
	defer mu.RUnlock()
	
	value, ok := s.data[key]
	return value, ok
}

func (s *Store) Delete(keys ...string) int {
	mu.Lock()
	defer mu.Unlock()
	count := 0

	for _, key := range keys {
		_, ok := s.data[key]
		if ok {
			delete(s.data, key)
			count++
		}
	}
	return count
}

func (s *Store) Exists(key string) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := s.data[key]
	return ok
}