package iterator

// Element represents a slice element
type Element any

// Elements represents a slice
type Elements []Element

// Iterator represents a iterator
type Iterator struct {
	el    Element
	elems []Element
	len   int
	pos   int
}

// New returns a new iterator
func New(elems []Element) *Iterator {
	return &Iterator{
		elems[0],
		elems,
		len(elems),
		0,
	}
}

// Next sets the next element
func (it *Iterator) Next() bool {
	if it.pos >= it.len-1 {
		return false
	}

	it.el = it.elems[it.pos]
	it.pos++

	return true
}

// Prev sets the previous element
func (it *Iterator) Prev() bool {
	if it.pos < 0 {
		return false
	}

	it.el = it.elems[it.pos]
	it.pos--

	return true
}

// Reset resets the iterator
func (it *Iterator) Reset() {
	it.pos = 0
	it.el = it.elems[0]
}

// CurrPos returns the current position
func (it *Iterator) CurrPos() int {
	return it.pos
}

// CurrElem returns the current element
func (it *Iterator) CurrElem() Element {
	return it.el
}

// Len returns the length
func (it *Iterator) Len() int {
	return it.len
}

// SetPos sets the position
func (it *Iterator) SetPos(pos int) bool {
	if it.len > pos {
		it.pos = pos
		it.el = it.elems[pos]
	}
	return it.len > pos
}
