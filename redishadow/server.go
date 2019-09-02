package main

// Server interface
type Server interface {
	Start()
	Stop()
	SetStorage(storage Storage)
}
