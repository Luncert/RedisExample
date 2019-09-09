package datastruct

import (
	"fmt"
	"testing"
)

func compareMethod(v1, v2 interface{}) (ret int) {
	a := v1.(int)
	b := v2.(int)
	if a > b {
		ret = 1
	} else if a == b {
		ret = 0
	} else {
		ret = -1
	}
	return
}

func createHeap() *BiHeap {
	heap := NewBiHeap(true, compareMethod)
	for i := 10; i > 0; i-- {
		heap.Append(i)
	}
	return heap
}

func TestBiHeap_Size(t *testing.T) {
	heap := createHeap()
	sz := heap.Size()
	if sz != 10 {
		t.Errorf("expect: %d, got: %d", 10, sz)
	}
}

func TestBiHeap_IsEmpty(t *testing.T) {
	heap := createHeap()
	if heap.IsEmpty() {
		t.Error("heap is empty")
	}
}

func TestBiHeap_Get(t *testing.T) {
	heap := NewBiHeap(true, compareMethod)
	heap.Append(10)
	v := heap.Get(0)
	if v != 10 {
		t.Errorf("expect: %d, got: %d", 10, v)
	}
}

// index out of range
func TestBiHeap_Get2(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			if err != "[BiHeap] index out of range" {
				t.Errorf("unexpected error: %v", err)
			}
		} else {
			t.Error("no error occurred")
		}
	}()
	heap := createHeap()
	heap.Get(10)
}

func TestBiHeap_Index(t *testing.T) {
	heap := createHeap()
	idx := heap.Index(10)
	if idx != 7 {
		t.Errorf("expect: %d, got: %d", 7, idx)
	}
}

func ExampleBiHeap_ForEach() {
	heap := createHeap()
	heap.ForEach(func(v interface{}) {
		fmt.Print(v, ",")
	})
	// Output: 1,2,5,4,3,9,6,10,7,8,
}

func ExampleBiHeap_Append() {
	heap := createHeap()
	fmt.Println(heap.data)
	// Output: [1 2 5 4 3 9 6 10 7 8]
}

func TestBiHeap_Pop(t *testing.T) {
	heap := createHeap()
	for i := 1; i <= 10; i++ {
		v := heap.Pop().(int)
		if v != i {
			t.Errorf("expect %d, got: %d", i, v)
		}
	}
}

func TestBiHeap_Set(t *testing.T) {
	heap := createHeap()
	// set 10 to -1
	heap.Set(7, -1)
	v := heap.Pop()
	if v != -1 {
		t.Errorf("expect: %d, got: %d", -1, v)
	}
}

func TestBiHeap_Set2(t *testing.T) {
	heap := createHeap()
	// set 1 to 11
	heap.Set(0, 11)
	v := heap.Pop()
	if v != 2 {
		t.Errorf("expect: %d, got: %d", 2, v)
	}
}

func TestBiHeap_Remove(t *testing.T) {
	heap := createHeap()
	r := []int{1, 2, 5, 4, 3, 9, 6, 10, 7, 8}
	for i := 9; i >= 0; i-- {
		v := heap.Remove(i)
		if v != r[i] {
			t.Errorf("expect: %d, got: %d", r[i], v)
		}
	}
}

func TestBiHeap_Reset(t *testing.T) {
	heap := createHeap()
	heap.Reset()
	if heap.size != 0 || len(heap.data) != 0 {
		t.Error("reset heap failed")
	}
}
