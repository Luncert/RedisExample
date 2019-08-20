package main

type MemoryStorage struct {
	data map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]string),
	}
}

func (s *MemoryStorage) SetString(k string, v string) {
	s.data[k] = v
}

func (s *MemoryStorage) GetString(k string) (v string, ok bool) {
	v, ok = s.data[k]
	return
}

func (s *MemoryStorage) DeleteKey(k string) {
	delete(s.data, k)
}
