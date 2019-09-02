package util

import (
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	s := make(chan bool, 1)
	timer := NewTimer(1*time.Second, func() {
		s <- true
	})
	timer.Start()
	<-s
	timer.Clear()
}
