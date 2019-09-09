package datastruct

type SkipList struct {
	head *skipListNode
	tail *skipListNode
	size int
}

func (s *SkipList) Size() int {
	return s.size
}

func (s *SkipList) IsEmpty() bool {
	return s.size == 0
}

func (s *SkipList) Get(index int) interface{} {
	return nil
}

func (s *SkipList) Index(value interface{}) int {
	return -1
}

func (s *SkipList) ForEach(call func(value interface{})) {

}

func (s *SkipList) Append(value interface{}) List {
	return s
}

func (s *SkipList) Pop() interface{} {
	return nil
}

func (s *SkipList) Set(index int, value interface{}) interface{} {
	return nil
}

func (s *SkipList) Remove(index int) interface{} {
	return nil
}

func (s *SkipList) Reset() {

}

func (s *SkipList) connect(node1, node2 *skipListNode, level int) {
}

type skipListNode struct {
	prev  []*skipListNode
	next  []*skipListNode
	value interface{}
}
