package main

// Storage ...
type Storage interface {
	Open()
	SetString(k string, v string)
	GetString(k string) (v string, ok bool)
	DeleteKey(k string)
	Close()
}
