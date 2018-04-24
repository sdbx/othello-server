package utils

import (
	"math/rand"
	"sync"
)

type keyStore struct {
	mu   sync.Mutex
	Keys map[string]bool
}

var letterRunes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var store = keyStore{
	Keys: make(map[string]bool),
}

func genKey(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (s *keyStore) gen() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	for {
		key := genKey(30)
		if !s.Keys[key] {
			s.Keys[key] = true
			return key
		}
	}
}

func GenKey() string {
	return store.gen()
}
