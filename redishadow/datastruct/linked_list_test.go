package datastruct

import (
	"fmt"
	"testing"
)

func TestAddNodeHead(t *testing.T) {
	list := NewLinkedList()

	node := list.AddNodeHead("head1")

	if list.head != node || node.value != "head1" || list.head.prev != nil || list.head.next != nil {
		t.Error("Incorrect head")
	}
	if list.tail != list.head {
		t.Error("Incorrect tail")
	}
	if list.len != 1 {
		t.Errorf("expect 1, got: %d\n", list.len)
	}

	node = list.AddNodeHead("head2")

	if list.head != node || node.value != "head2" || list.head.prev != nil || list.head.next != list.tail {
		t.Error("Incorrect head")
	}
	if list.tail == nil || list.tail.value != "head1" || list.tail.prev != list.head || list.tail.next != nil {
		t.Error("Incorrect tail")
	}
	if list.len != 2 {
		t.Errorf("expect 2, got: %d\n", list.len)
	}
}

func TestAddNodeTail(t *testing.T) {
	list := NewLinkedList()

	node := list.AddNodeTail("tail1")

	if list.head != node || node.value != "tail1" || list.head.prev != nil || list.head.next != nil {
		t.Error("Incorrect head")
	}
	if list.tail != list.head {
		t.Error("Incorrect tail")
	}
	if list.len != 1 {
		t.Errorf("expect 1, got: %d\n", list.len)
	}

	node = list.AddNodeTail("tail2")

	if list.head == nil || list.head.value != "tail1" || list.head.prev != nil || list.head.next != node {
		t.Error("Incorrect head")
	}
	if list.tail != node || node.value != "tail2" || list.tail.prev != list.head || list.tail.next != nil {
		t.Error("Incorrect tail")
	}
	if list.len != 2 {
		t.Errorf("expect 2, got: %d\n", list.len)
	}
}

func TestInsertBefore(t *testing.T) {
	list := NewLinkedList()
	node := &LinkedNode{value: "node"}
	list.head = node
	list.tail = node
	list.len = 1

	node = list.InsertBefore(0, "node1")
	if list.head != node || node.value != "node1" || list.head.prev != nil || list.head.next != list.tail {
		t.Error("Incorrect head")
	}
	if list.tail == nil || list.tail.value != "node" || list.tail.prev != node || list.tail.next != nil {
		t.Error("Incorrect tail")
	}
	if list.len != 2 {
		t.Errorf("expect 2, got: %d\n", list.len)
	}

	node = list.InsertBefore(0, "node2")
	if list.head != node || node.value != "node2" || list.head.prev != nil || list.head.next == nil {
		t.Error("Incorrect head")
	}
	if list.tail == nil || list.tail.value != "node" || list.tail.prev == nil || list.tail.next != nil {
		t.Error("Incorrect tail")
	}
	if list.len != 3 {
		t.Errorf("expect 3, got: %d\n", list.len)
	}
}

func TestInsertAfter(t *testing.T) {
	list := NewLinkedList()
	node := &LinkedNode{value: "node"}
	list.head = node
	list.tail = node
	list.len = 1

	list.InsertAfter(0, "node1")

	if list.head == nil || list.head.value != "node" {
		t.Error("Incorrect head")
	}
	if list.head.prev != nil || list.head.next != list.tail {
		t.Error("Incorrect head")
	}

	if list.tail == nil || list.tail.value != "node1" {
		t.Error("Incorrect tail")
	}
	if list.tail.next != nil || list.tail.prev != list.head {
		t.Error("Incorrect tail")
	}

	if list.len != 2 {
		t.Errorf("expect 2, got: %d\n", list.len)
	}
}

func TestAppendPop(t *testing.T) {
	list := NewLinkedList()
	list.Append("1")

	if list.head == nil || list.head.value != "1" || list.head.prev != nil || list.head.next != nil {
		t.Error("Incorrect head")
	}
	if list.tail != list.head {
		t.Error("Incorrect tail")
	}
	if list.len != 1 {
		t.Errorf("expect 1, got: %d\n", list.len)
	}

	list.Append("2")
	if list.head == nil || list.head.value != "1" || list.head.prev != nil || list.head.next != list.tail {
		t.Error("Incorrect head")
	}
	if list.tail == nil || list.tail.value != "2" || list.tail.next != nil || list.tail.prev != list.head {
		t.Error("Incorrect tail")
	}
	if list.len != 2 {
		t.Errorf("expect 1, got: %d\n", list.len)
	}

	v := list.Pop().(string)
	if v != "2" {
		t.Errorf("expect 2, got: %s", v)
	}
	if list.head == nil || list.head.value != "1" || list.head.prev != nil || list.head.next != nil {
		t.Error("Incorrect head")
	}
	if list.tail != list.head {
		t.Error("Incorrect tail")
	}
	if list.len != 1 {
		t.Errorf("expect 1, got: %d\n", list.len)
	}

	v = list.Pop().(string)
	if v != "1" {
		t.Errorf("expect 1, got: %s", v)
	}
	if list.head != nil {
		t.Error("Incorrect head")
	}
	if list.tail != nil {
		t.Error("Incorrect tail")
	}
	if list.len != 0 {
		t.Errorf("expect 0, got: %d\n", list.len)
	}
}

func TestFirstLast(t *testing.T) {
	list := NewLinkedList()
	if ret := list.First(); ret != nil {
		t.Errorf("expect nil, got: %v", ret)
	}
	if ret := list.Last(); ret != nil {
		t.Errorf("expect nil, got: %v", ret)
	}

	list.Append("1")
	if ret := list.First(); ret == nil {
		t.Errorf("epxect *LinkedNode, got nil")
	} else if ret.value != "1" {
		t.Errorf("expect 1, got: %s", ret.value)
	}
	if ret := list.Last(); ret == nil {
		t.Errorf("epxect *LinkedNode, got nil")
	} else if ret.value != "1" {
		t.Errorf("expect 1, got: %s", ret.value)
	}
}

func TestAtSearchIndex(t *testing.T) {
	list := NewLinkedList().
		Append("node1").
		Append("node2").
		Append("node3").(*LinkedList)

	list.SetMatchMethod(func(v1, v2 interface{}) bool {
		return v1 == v2
	})

	node1 := list.At(0)
	node2 := list.At(1)
	node3 := list.At(2)

	if node1 == nil || node1.value != "node1" || node1.prev != nil || node1.next != node2 {
		t.Error("Incorrect node1")
	}
	if node2 == nil || node2.value != "node2" || node2.prev != node1 || node2.next != node3 {
		t.Error("Incorrect node2")
	}
	if node3 == nil || node3.value != "node3" || node3.prev != node2 || node3.next != nil {
		t.Error("Incorrect node3")
	}

	tmp := list.Search("node2")
	if tmp != node2 {
		t.Error("Incorrect result")
	}

	i := list.Index("node2")
	if i != 1 {
		t.Error("Incorrect result")
	}
}

func ExampleForEach() {
	list := NewLinkedList().
		Append("testing").
		Append("ForEach").
		Append("function").
		Append("#1")
	list.ForEach(func(value interface{}) {
		fmt.Printf("%v,", value)
	})
	// Output: testing,ForEach,function,#1,
}

func TestRemoveAt(t *testing.T) {
	list := NewLinkedList().
		Append("node1").
		Append("node2").
		Append("node3")

	v := list.Remove(2)
	if v != "node3" {
		t.Error("Incorrect result")
	}
	v = list.Remove(1)
	if v != "node2" {
		t.Error("Incorrect result")
	}
	v = list.Remove(0)
	if v != "node1" {
		t.Error("Incorrect result")
	}
	if !list.IsEmpty() {
		t.Error("Incorrect result")
	}
}

func TestReset(t *testing.T) {
	list := NewLinkedList().
		Append("node1").
		Append("node2").
		Append("node3").(*LinkedList)

	list.Reset()
	if list.head != nil || list.tail != nil || list.len != 0 {
		t.Error("Incorrect result")
	}
}
