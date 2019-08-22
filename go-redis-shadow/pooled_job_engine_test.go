package main

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkPooledJobEngine(b *testing.B) {
	jobEngine := NewPooledJobEngine(512, 32, 128, 5*time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := jobEngine.Execute(taskCallable, i); err != nil {
			fmt.Println(err)
			b.Fail()
		}
	}
	jobEngine.Shutdown()
}

func taskCallable(args ...interface{}) {
}
