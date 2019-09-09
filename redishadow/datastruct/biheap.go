package datastruct

// BiHeap supports both MinHeap and MaxHeap,
// but its elements are out of order.
type BiHeap struct {
	minHeap       bool
	compareMethod func(e1, e2 interface{}) int
	size          int
	data          []interface{}
}

func NewBiHeap(minHeap bool,
	compareMethod func(e1, e2 interface{}) int) *BiHeap {
	return &BiHeap{
		minHeap:       minHeap,
		compareMethod: compareMethod,
		data:          make([]interface{}, 0),
	}
}

func (b *BiHeap) Size() int {
	return b.size
}

func (b *BiHeap) IsEmpty() bool {
	return b.size == 0
}

// Get, this function doesn't make sense because the elements of BiHeap
// are out of order.
func (b *BiHeap) Get(index int) interface{} {
	b.checkIndex(index)
	return b.data[index]
}

func (b *BiHeap) Index(value interface{}) int {
	return b.index(0, value)
}

func (b *BiHeap) index(baseNode int, value interface{}) int {
	idx := -1
	if baseNode < b.size {
		v := b.data[baseNode]
		r := b.compareMethod(v, value)
		if r == 0 {
			idx = baseNode
		} else if b.minHeap && r < 0 || !b.minHeap && r > 0 {
			idx = b.index(baseNode*2+1, value)
			if idx == -1 {
				idx = b.index(baseNode*2+2, value)
			}
		}
	}
	return idx
}

// ForEach, honestly, this function doesn't make sense.
func (b *BiHeap) ForEach(call func(value interface{})) {
	for _, v := range b.data {
		call(v)
	}
}

func (b *BiHeap) Append(value interface{}) List {
	b.data = append(b.data, value)
	b.percolateUp(b.size)
	b.size++
	return b
}

/*
[Important] Odd index refers to left children, even index refers to right children.
*/

func (b *BiHeap) percolateUp(idx int) {
	// If idx value is odd, then the result of idx / 2 and (idx -1) /2 are equals,
	// because this is integer division, e.g. idx = 3, then parentIdx should be 1.
	// If idx value if even, we should calculate (idx - 1) / 2 to get its parent idx,
	// e.g. idx = 6, (idx - 1) / 2 = 4, but idx / 2 = 3, 3 is not correct.
	parentIdx := (idx - 1) / 2
	// The element with index 0, is the root element of this binary tree.
	// In other words, other elements all have a parent.
	for idx > 0 && b.percolatePart(parentIdx, idx) {
		idx = parentIdx
		parentIdx = (idx - 1) / 2
	}
}

func (b *BiHeap) percolateDown(idx int) {
	lcIdx := idx*2 + 1
	rcIdx := lcIdx + 1
	if lcIdx < b.size && b.percolatePart(idx, lcIdx) {
		b.percolateDown(lcIdx)
	}
	if rcIdx < b.size && b.percolatePart(idx, rcIdx) {
		b.percolateDown(rcIdx)
	}
}

func (b *BiHeap) percolatePart(pIdx, cIdx int) (ok bool) {
	r := b.compareMethod(b.data[pIdx], b.data[cIdx])
	if b.minHeap && r > 0 || !b.minHeap && r < 0 {
		tmp := b.data[pIdx]
		b.data[pIdx] = b.data[cIdx]
		b.data[cIdx] = tmp
		ok = true
	}
	return
}

// Pop: delete and return the first element of BiHeap,
// for MinHeap it's the minimum value, for MaxHeap it's the maximum value.
func (b *BiHeap) Pop() interface{} {
	return b.Remove(0)
}

func (b *BiHeap) Set(index int, value interface{}) interface{} {
	b.checkIndex(index)
	oldValue := b.data[index]
	r := b.compareMethod(oldValue, value)
	b.data[index] = value
	if r < 0 {
		// new value is bigger
		b.percolateDown(index)
	} else if r > 0 {
		// new value if smaller
		b.percolateUp(index)
	}
	return oldValue
}

func (b *BiHeap) Remove(index int) interface{} {
	b.checkIndex(index)
	// Move the last element to the position of deleted element,
	// and then perform percolateDown.
	value := b.data[index]
	b.size--
	b.data[index] = b.data[b.size]
	b.data[b.size] = nil
	b.percolateDown(index)

	// release unused space
	if len(b.data)/2 >= b.size {
		b.data = b.data[:b.size]
	}
	return value
}

func (b *BiHeap) Reset() {
	b.size = 0
	b.data = make([]interface{}, 0)
}

func (b *BiHeap) checkIndex(index int) {
	if index < 0 || index >= b.size {
		panic("[BiHeap] index out of range")
	}
}
