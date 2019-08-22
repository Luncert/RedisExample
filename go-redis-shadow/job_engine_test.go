package main

import (
	"testing"
	"time"
	"fmt"
)

func taskCallable(args ...interface{}) {
}

func BenchmarkJobEngine(b *testing.B) {
	jobEngine := NewJobEngine(128, 32)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := jobEngine.Execute(taskCallable, i); err != nil {
			fmt.Println(err)
			b.Fail()
		}
	}
	jobEngine.Shutdown()
}

func TestJobEngineX(t *testing.T) {
	jobEngine := NewJobEngineX(128, 8, 16, 5 * time.Second)
	if err := jobEngine.Execute(taskCallable, "test"); err != nil {
		fmt.Println(err)
		t.Fail()
	}
	jobEngine.Shutdown()
}

func BenchmarkJobEngineX(b *testing.B) {
	jobEngine := NewJobEngineX(128, 8, 32, 5 * time.Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := jobEngine.Execute(taskCallable, i); err != nil {
			fmt.Println(err)
			b.Fail()
		}
	}
	jobEngine.Shutdown()
}