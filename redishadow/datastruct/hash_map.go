package datastruct

const (
	defaultCapacity = 8
)

type HashMap struct {
	size  int
	cap   int
	table []tableEntry
}

func (h *HashMap) indexEntry(key interface{}) tableEntry {
	return h.table[h.hash(key)&(h.cap-1)]
}

func (h *HashMap) hash(key interface{}) int {
	return 0
}

func (h *HashMap) Size() int {
	return h.size
}

func (h *HashMap) IsEmpty() bool {
	return h.size == 0
}

func (h *HashMap) HasKey(key interface{}) bool {
	entry := h.indexEntry(key)
	_, ok := entry.search(key)
	return ok
}

func (h *HashMap) HasValue(value interface{}) bool {
	for _, entry := range h.table {
		if entry != nil && entry.hasValue(value) {
			return true
		}
	}
	return false
}

func (h *HashMap) Keys() []interface{} {
	keys := make([]interface{}, 0)
	for _, entry := range h.table {
		if entry != nil {
			keys = append(keys, entry.keys()...)
		}
	}
	return keys
}

func (h *HashMap) Values() []interface{} {
	values := make([]interface{}, 0)
	for _, entry := range h.table {
		if entry != nil {
			values = append(values, entry.values()...)
		}
	}
	return values
}

func (h *HashMap) Get(key interface{}) (interface{}, bool) {
	entry := h.indexEntry(key)
	return entry.search(key)
}

func (h *HashMap) Put(key interface{}, value interface{}) {
	entry := h.indexEntry(key)
	entry.add(key, value)
}

func (h *HashMap) PutAll(m HashMap) {
	m.ForEach(func(key interface{}, value interface{}) {
		h.Put(key, value)
	})
}

func (h *HashMap) Remove(key interface{}) interface{} {
	entry := h.indexEntry(key)
	return entry.remove(key)
}

func (h *HashMap) Clear() {
	h.size = 0
	h.cap = defaultCapacity
	h.table = make([]tableEntry, defaultCapacity)
}

func (h *HashMap) ForEach(call func(key interface{}, value interface{})) {
	for _, entry := range h.table {
		if entry != nil {
			entry.forEach(call)
		}
	}
}

type tableEntry interface {
	size() int
	add(key interface{}, value interface{})
	search(key interface{}) (value interface{}, ok bool)
	hasValue(value interface{}) bool
	remove(key interface{}) (value interface{})
	keys() []interface{}
	values() []interface{}
	forEach(call func(key interface{}, value interface{}))
}

type keyValue struct {
	key   interface{}
	value interface{}
}

type linkedListTableEntry struct {
	list *LinkedList
}

func newLinkedListTableEntry() *linkedListTableEntry {
	return &linkedListTableEntry{list: NewLinkedList()}
}

func (l *linkedListTableEntry) size() int {
	return l.list.Size()
}

func (l *linkedListTableEntry) add(key interface{}, value interface{}) {
	l.list.AddNodeHead(keyValue{key: key, value: value})
}

func (l *linkedListTableEntry) search(key interface{}) (value interface{}, ok bool) {
	return nil, false
}

func (l *linkedListTableEntry) hasValue(value interface{}) bool {
	return false
}

func (l *linkedListTableEntry) remove(key interface{}) (value interface{}) {
	return nil
}

func (l *linkedListTableEntry) keys() []interface{} {
	return nil
}

func (l *linkedListTableEntry) values() []interface{} {
	return nil
}

func (l *linkedListTableEntry) forEach(call func(key interface{}, value interface{})) {
}
