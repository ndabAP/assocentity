package iterator

// Iterator represents a iterator
type Iterator struct {
	slice *[]string
	pos   int
	el    string
	init  bool
}

// New returns a new iterator
func New(slice *[]string) *Iterator {
	return &Iterator{
		slice,
		0,
		(*slice)[0],
		true,
	}
}

// Next sets the next element
func (g *Iterator) Next() bool {
	if g.pos+1 > g.Len()-1 {
		return false
	}

	if g.init {
		g.init = false

		return true
	}

	g.pos++
	g.el = (*g.slice)[g.pos]

	return true
}

// Prev sets the next element
func (g *Iterator) Prev() bool {
	if g.pos-1 < 0 {
		return false
	}

	if g.init {
		g.init = false

		return true
	}

	g.pos--
	g.el = (*g.slice)[g.pos]

	return true
}

// Reset resets the iterator
func (g *Iterator) Reset() {
	g.pos = 0
	g.el = (*g.slice)[0]
	g.init = true
}

// CurrPos returns the current position
func (g *Iterator) CurrPos() int {
	return g.pos
}

// CurrElem returns the current element
func (g *Iterator) CurrElem() string {
	return g.el
}

// Len returns the length
func (g *Iterator) Len() int {
	return len((*g.slice))
}

// SetPos sets the position
func (g *Iterator) SetPos(pos int) bool {
	if g.Len() > pos {
		g.pos = pos
		g.el = (*g.slice)[pos]
	}

	return g.Len() > pos
}
