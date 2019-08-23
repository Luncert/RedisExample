package main

import (
	"fmt"
	"testing"
)

func BenchmarkJobEngine(b *testing.B) {
	taskCallable := func(args ...interface{}) {}
	jobEngine := NewJobEngine(256, 128)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := jobEngine.Execute(taskCallable, i); err != nil {
			fmt.Println(err)
			b.Fail()
		}
	}
	jobEngine.Shutdown()
}
