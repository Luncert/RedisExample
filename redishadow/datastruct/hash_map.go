package datastruct

type HashMap struct {
	table []tableEntry
}

func (h *HashMap) Size() int {
	return 0
}

func (h *HashMap) IsEmpty() bool {
	return false
}

func (h *HashMap) HasKey(key interface{}) bool {
	return false
}

func (h *HashMap) HasValue(value interface{}) bool {
	return false
}

func (h *HashMap) Keys() []interface{} {
	return nil
}

func (h *HashMap) Values() []interface{} {
	return nil
}

func (h *HashMap) Get(key interface{}) (interface{}, bool) {
	return nil, false
}

func (h *HashMap) Put(key interface{}, value interface{}) {

}

func (h *HashMap) PutAll(m HashMap) {

}

func (h *HashMap) Remove(key interface{}) {

}

func (h *HashMap) Clear() {

}

type tableEntry interface {
	size() int
	add(key interface{}, value interface{})
	search(key interface{}) (value interface{})
	hasValue(value interface{}) bool
	remove(key interface{}) (value interface{})
}
