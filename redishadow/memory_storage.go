package main

// MemoryStorage ...
type MemoryStorage struct {
	data map[string]string
}

// NewMemoryStorage ...
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]string),
	}
}

func (s *MemoryStorage) Open() {

}

// SetString ...
func (s *MemoryStorage) SetString(k string, v string) {
	s.data[k] = v
}

// GetString ...
func (s *MemoryStorage) GetString(k string) (v string, ok bool) {
	v, ok = s.data[k]
	return
}

// DeleteKey ...
func (s *MemoryStorage) DeleteKey(k string) {
	delete(s.data, k)
}

func (s *MemoryStorage) Close() {

}
