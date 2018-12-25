package assocentity

import (
	"testing"
)

var (
	text      = []byte("Hello, my name is John Max. I'm the best human and I'm John. The real John Max, oh yes!")
	blacklist = []string{"while", "is", "yes", ",", ".", "!", "?", ":", ";"}
	whitelist = []string{"nice"}
	entity    = "John Max"
	separator = []byte(" ")
)

func TestNewFilter(t *testing.T) {
	f := NewFilter(blacklist, whitelist, entity, separator)

	if whitelist := len(f.whitelist); whitelist != 1 {
		t.Errorf("Filter blacklist incorrect, got %d., want: %d.", whitelist, 1)
	}
	if blacklist := len(f.blacklist); blacklist != 9 {
		t.Errorf("Filter blacklist incorrect, got %d., want: %d.", blacklist, 9)
	}
	if f.entity != entity {
		t.Errorf("Filter entities incorrect, got %s., want: %d.", f.entity, 1)
	}
}

func TestTraversableLatin(t *testing.T) {
	f := NewFilter(blacklist, whitelist, entity, separator)

	traversable := TraversableLatin(text, f)

	if len(traversable) != 17 {
		t.Errorf("Traversable latin incorrect, got %d., want: %d.", len(traversable), 17)
	}
}

func TestGraph(t *testing.T) {
	f := NewFilter(blacklist, whitelist, entity, separator)

	traversable := TraversableLatin(text, f)

	Graph(f, traversable)
}
