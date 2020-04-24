package app

import (
	"encoding/json"
	"sync"
)

type jsStorage struct {
	name  string
	key   []byte
	mutex sync.RWMutex
}

func newJSStorage(name string) *jsStorage {
	u := Window().URL()

	key := []byte(u.Scheme + "(*_*)" + u.Host)
	for len(key) < 32 {
		key = append(key, 'o')
	}
	key = key[:32]

	return &jsStorage{
		name: name,
		key:  key,
	}
}

func (s *jsStorage) Set(k string, v interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	Window().Get(s.name).Call("setItem", k, string(b))
	return nil
}

func (s *jsStorage) Get(k string, v interface{}) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	item := Window().Get(s.name).Call("getItem", k)
	if !item.Truthy() {
		return nil
	}

	return json.Unmarshal([]byte(item.String()), v)
}

func (s *jsStorage) Del(k string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	Window().Get(s.name).Call("removeItem", k)
}

func (s *jsStorage) Clear() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	Window().Get(s.name).Call("clear")
}
