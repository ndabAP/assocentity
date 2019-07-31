package generator

// Generator represents a generator
type Generator struct {
	slice []string
	pos   int
	el    string
	init  bool
}

// New returns a new generator
func New(slice []string) *Generator {
	return &Generator{
		slice,
		0,
		slice[0],
		true,
	}
}

// Next sets the next element
func (g *Generator) Next() bool {
	if g.pos+1 > len(g.slice)-1 {
		return false
	}

	if g.init {
		g.init = false

		return true
	}

	g.pos++
	g.el = g.slice[g.pos]

	return true
}

// Prev sets the next element
func (g *Generator) Prev() bool {
	if g.pos-1 < 0 {
		return false
	}

	if g.init {
		g.init = false

		return true
	}

	g.pos--
	g.el = g.slice[g.pos]

	return true
}

// Reset resets the generator
func (g *Generator) Reset() {
	g.pos = 0
	g.el = g.slice[0]
	g.init = true
}

// CurrPos returns the current position
func (g *Generator) CurrPos() int {
	return g.pos
}

// CurrElem returns the current element
func (g *Generator) CurrElem() string {
	return g.el
}

// Len returns the length
func (g *Generator) Len() int {
	return len(g.slice)
}

// SetPos sets the position
func (g *Generator) SetPos(pos int) bool {
	if len(g.slice) > pos {
		g.pos = pos
		g.el = g.slice[pos]
		g.init = true
	}

	return len(g.slice) > pos
}
