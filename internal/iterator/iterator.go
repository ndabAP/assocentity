package iterator

// Iterator represents a iterator
type Iterator[T any] struct {
	el    T
	elems []T
	len   int
	pos   int

	init bool
}

// New returns a new iterator
func New[T any](elems []T) *Iterator[T] {
	return &Iterator[T]{
		elems[0],
		elems,
		len(elems),
		0,
		true,
	}
}

// Next sets the next element
func (it *Iterator[T]) Next() bool {
	// Delays the index
	if it.init {
		it.init = false
		return true
	}

	// We increment before assigning since we used "init"
	it.pos++
	if it.pos >= it.len {
		return false
	}
	it.el = it.elems[it.pos]
	return true
}

// Prev sets the previous element
func (it *Iterator[T]) Prev() bool {
	if it.init {
		it.el = it.elems[0]
		it.init = false
		return true
	}

	it.pos--
	if it.pos < 0 {
		return false
	}
	it.el = it.elems[it.pos]
	return true
}

func (it *Iterator[T]) Elems() []T {
	return it.elems
}

// Reset resets the iterator
func (it *Iterator[T]) Reset() {
	it.pos = 0
	it.el = it.elems[0]
	it.init = true
}

// CurrPos returns the current position
func (it *Iterator[T]) CurrPos() int {
	return it.pos
}

// CurrElem returns the current element
func (it *Iterator[T]) CurrElem() T {
	return it.el
}

// Len returns the length
func (it *Iterator[T]) Len() int {
	return it.len
}

// SetPos sets the position
func (it *Iterator[T]) SetPos(pos int) *Iterator[T] {
	if it.len > pos && pos >= 0 {
		it.pos = pos
		it.el = it.elems[it.pos]
	}
	return it
}
