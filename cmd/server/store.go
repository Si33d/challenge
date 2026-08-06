package main

import (
	"sync"
	"time"
)

type Store struct {
	data map[string]string
	expiry map[string]time.Time
}

var mu sync.RWMutex

var store = NewStore() //level3 done here only

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
		expiry: make(map[string]time.Time),
	}
}

func (s *Store) Set(key, value string,expiry *time.Time) {
	mu.Lock()
	defer mu.Unlock()
	s.data[key] = value

	if expiry!=nil{
		s.expiry[key]=*expiry
	}else{
		delete(s.expiry,key)
	}
}

func (s *Store) Get(key string) (string, bool) {
	mu.RLock()
	defer mu.RUnlock()

	value, ok := s.data[key]
	if !ok{
		return "",false
	}
	exptime,isexpire:=s.expiry[key]
	if isexpire{
		if time.Now().After(exptime){
			delete(s.data,key)
			delete(s.expiry,key)

			return "",false
		}
	}
	return value,true
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

func(s *Store)Expire(key string,seconds int)bool{
	mu.Lock()
	defer mu.Unlock()

	_,ok:=s.data[key]
	if !ok{
		return false
	}
	s.expiry[key]=time.Now().Add(time.Duration(seconds)*time.Second)
	return true
}

func(s *Store)TTL(key string)int{
	mu.Lock()
	defer mu.Unlock()

	_,ok:=s.data[key]
	if !ok{
		return -2
	}
	exptime,hasexpire:=s.expiry[key]
	if !hasexpire{
		return -1
	}
	if time.Now().After(exptime){
		delete(s.data,key)
		delete(s.expiry,key)
		return -2
	}
	return int(time.Until(exptime).Seconds())
}