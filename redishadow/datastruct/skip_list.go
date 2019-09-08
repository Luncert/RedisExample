package datastruct

type SkipList struct {
	head *skipListNode
	tail *skipListNode
	len  int
}

func (s *SkipList) Size() int {

}

func (s *SkipList) IsEmpty() bool {

}

func (s *SkipList) Get(index int) (interface{}, bool) {

}

func (s *SkipList) Index(value interface{}) int {

}

func (s *SkipList) ForEach(call func(value interface{})) {

}

func (s *SkipList) Append(value interface{}) List {

}

func (s *SkipList) Pop() interface{} {

}

func (s *SkipList) Set(index int, value interface{}) interface{} {

}

func (s *SkipList) Remove(index int) interface{} {

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
