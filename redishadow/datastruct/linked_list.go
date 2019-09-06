package datastruct

type LinkedList struct {
	head  *LinkedNode
	tail  *LinkedNode
	len   int
	match func(v1, v2 interface{}) bool
}

func NewLinkedList() *LinkedList {
	return &LinkedList{
		head:  nil,
		tail:  nil,
		len:   0,
		match: nil,
	}
}

func (l *LinkedList) SetMatchMethod(match func(v1, v2 interface{}) bool) {
	l.match = match
}

func (l *LinkedList) Size() int {
	return l.len
}

func (l *LinkedList) IsEmpty() bool {
	return l.len == 0
}

// get

func (l *LinkedList) First() *LinkedNode {
	return l.head
}

func (l *LinkedList) At(index int) *LinkedNode {
	l.checkIndex(index)
	node := l.head
	for i := 0; i != index; i++ {
		node = node.next
	}
	return node
}

func (l *LinkedList) Last() *LinkedNode {
	return l.tail
}

func (l *LinkedList) Search(value interface{}) *LinkedNode {
	l.checkMatchMethod()
	node := l.head
	for i := 0; i < l.len; i++ {
		if l.match(node.value, value) {
			return node
		}
		node = node.next
	}
	return nil
}

func (l *LinkedList) Index(value interface{}) int {
	l.checkMatchMethod()
	node := l.head
	for i := 0; i < l.len; i++ {
		if l.match(node.value, value) {
			return i
		}
		node = node.next
	}
	return -1
}

// create

func (l *LinkedList) AddNodeHead(value interface{}) *LinkedNode {
	node := &LinkedNode{
		prev:  nil,
		next:  l.head,
		value: value,
	}
	if l.head != nil {
		l.head.prev = node
	} else {
		l.tail = node
	}
	l.head = node

	l.len++
	return node
}

func (l *LinkedList) AddNodeTail(value interface{}) *LinkedNode {
	node := &LinkedNode{
		prev:  l.tail,
		next:  nil,
		value: value,
	}
	if l.tail != nil {
		l.tail.next = node
	} else {
		l.head = node
	}
	l.tail = node

	l.len++
	return node
}

func (l *LinkedList) InsertBefore(index int, value interface{}) *LinkedNode {
	nextNode := l.At(index)
	node := &LinkedNode{
		prev:  nextNode.prev,
		next:  nextNode,
		value: value,
	}
	if nextNode == l.head {
		l.head = node
	} else {
		node.prev.next = node
	}
	nextNode.prev = node

	l.len++
	return node
}

func (l *LinkedList) InsertAfter(index int, value interface{}) *LinkedNode {
	prevNode := l.At(index)
	node := &LinkedNode{
		prev:  prevNode,
		next:  prevNode.next,
		value: value,
	}
	if prevNode == l.tail {
		l.tail = node
	} else {
		node.next.prev = node
	}
	prevNode.next = node

	l.len++
	return node
}

func (l *LinkedList) Append(value interface{}) *LinkedList {
	node := &LinkedNode{
		prev:  l.tail,
		next:  nil,
		value: value,
	}
	if l.len == 0 {
		l.head = node
	} else {
		l.tail.next = node
	}
	l.tail = node
	l.len++
	return l
}

func (l *LinkedList) Pop() interface{} {
	if l.IsEmpty() {
		panic("LinkedList: pop operation on empty collection")
	}
	node := l.tail
	l.tail = node.prev
	if l.tail == nil {
		l.head = nil
	} else {
		l.tail.next = nil
	}
	l.len--
	return node.value
}

// remove

func (l *LinkedList) RemoveAt(index int) interface{} {
	l.checkIndex(index)

	var node *LinkedNode
	if l.len == 1 {
		node = l.head
		l.head = nil
		l.tail = nil
		l.len = 0
	} else {
		node = l.head
		for i := 0; i != index; i++ {
			node = node.next
		}

		if node == l.head {
			l.head = node.next
			l.head.prev = nil
		} else if node == l.tail {
			l.tail = node.prev
			l.tail.next = nil
		} else {
			prevNode := node.prev
			nextNode := node.next
			prevNode.next = nextNode
			nextNode.prev = prevNode
		}

		l.len--
	}
	return node.value
}

func (l *LinkedList) Reset() {
	l.head = nil
	l.tail = nil
	l.len = 0
}

func (l *LinkedList) checkMatchMethod() {
	if l.match == nil {
		panic("LinkedList: no implementation of match method")
	}
}

func (l *LinkedList) checkIndex(index int) {
	if index < 0 || index >= l.len {
		panic("LinkedList: index out of range")
	}
}

type LinkedNode struct {
	prev  *LinkedNode
	next  *LinkedNode
	value interface{}
}

func (l *LinkedNode) PrevNode() *LinkedNode {
	return l.prev
}

func (l *LinkedNode) NextNode() *LinkedNode {
	return l.next
}

func (l *LinkedNode) SetValue(newValue interface{}) interface{} {
	defer func() {
		l.value = newValue
	}()
	return l.value
}

func (l *LinkedNode) Value() interface{} {
	return l.value
}
