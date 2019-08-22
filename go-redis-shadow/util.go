package main

import (
	"time"
)

// Timer ...
type Timer struct {
	stopSignal chan bool
	timeout    time.Duration
	callback   func()
	stopped    bool
}

// NewTimer ...
func NewTimer(timeout time.Duration, callback func()) *Timer {
	return &Timer{
		stopSignal: make(chan bool, 1),
		timeout:    timeout,
		callback:   callback,
		stopped:    false,
	}
}

// Start ...
func (t *Timer) Start() {
	t.stopped = false
	go func() {
		select {
		case <-time.After(t.timeout):
			t.callback()
		case <-t.stopSignal:
		}
		t.stopped = true
	}()
}

func (t *Timer) Reset() {
	if !t.stopped {
		t.Clear()
	}
	t.Start()
}

// Clear ...
func (t *Timer) Clear() {
	t.stopSignal <- true
}
