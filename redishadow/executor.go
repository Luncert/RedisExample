package main

// Executor ...
type Executor interface {
	Execute(input []byte, s *Storage) []byte
}
