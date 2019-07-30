package generator

// Generator represents a generator
type Generator struct {
	slice []string
	pos   int
	el    string
}

// New returns a new generator
func New(slice []string) *Generator {
	return &Generator{
		slice,
		0,
		slice[0],
	}
}

// Next sets the next element
func (g *Generator) Next() bool {
	if g.pos+1 > len(g.slice) {
		return false
	}

	if g.pos == 0 {
		if g.pos+1 > len(g.slice) {
			return false
		}

		g.pos++

		return true
	}

	g.pos++
	g.el = g.slice[g.pos]

	return true
}

// Prev sets the next element
func (g *Generator) Prev() bool {
	if g.pos == 0 {
		return false
	}

	g.pos--
	g.el = g.slice[g.pos]

	return true
}

// Reset resets the slice
func (g *Generator) Reset() {
	g.pos = 0
	g.el = g.slice[0]
}

// GetCurrPos returns the current position
func (g *Generator) GetCurrPos() int {
	return g.pos
}

// GetCurrElem returns the current element
func (g *Generator) GetCurrElem() string {
	return g.el
}
