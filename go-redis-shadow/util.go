package main

import (
	"time"
)

// Timer ...
type Timer struct {
	stopSignal chan bool
	timeout    time.Duration
	callback   func()
}

// NewTimer ...
func NewTimer(timeout time.Duration, callback func()) *Timer {
	return &Timer{
		stopSignal: make(chan bool, 1),
		timeout:    timeout,
		callback:   callback,
	}
}

// Start ...
func (t *Timer) Start() {
	go func() {
		select {
		case <-time.After(t.timeout):
			t.callback()
		case <-t.stopSignal:
		}
	}()
}

// Clear ...
func (t *Timer) Clear() {
	t.stopSignal <- true
}
