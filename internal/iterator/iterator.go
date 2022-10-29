package iterator

// Iterator represents an iterator
type Iterator[T any] struct {
	el    T
	elems []T
	len   int
	pos   int
}

// New returns a new iterator
func New[T any](elems []T) *Iterator[T] {
	return &Iterator[T]{
		elems[0],
		elems,
		len(elems),
		-1,
	}
}

// Next sets the next element
func (it *Iterator[T]) Next() bool {
	// We increment before assigning since we used "init"
	if it.pos+1 >= it.len {
		return false
	}

	it.pos++
	it.el = it.elems[it.pos]
	return true
}

// Prev sets the previous element
func (it *Iterator[T]) Prev() bool {
	if it.pos-1 < 0 {
		return false
	}

	it.pos--
	it.el = it.elems[it.pos]
	return true
}

func (it *Iterator[T]) Elems() []T {
	return it.elems
}

// Reset resets the iterator
func (it *Iterator[T]) Reset() *Iterator[T] {
	it.pos = -1
	it.el = it.elems[0]
	return it
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
	it.pos = pos
	it.setEl()
	return it
}

func (it *Iterator[T]) Rewind(pos int) *Iterator[T] {
	it.pos -= pos
	it.setEl()
	return it
}

func (it *Iterator[T]) Foward(pos int) *Iterator[T] {
	it.pos += pos
	it.setEl()
	return it
}

func (it *Iterator[T]) setEl() {
	if len(it.elems)-1 > it.pos && it.pos >= 0 {
		it.el = it.elems[it.pos]
	}
}

// TODO: Save(), Revert(). Revert iterator
