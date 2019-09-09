package datastruct

type List interface {
	Size() int
	IsEmpty() bool
	Get(index int) interface{}
	Index(value interface{}) int
	ForEach(call func(value interface{}))
	Append(value interface{}) List
	Pop() interface{}
	Set(index int, value interface{}) interface{}
	Remove(index int) interface{}
	Reset()
}
