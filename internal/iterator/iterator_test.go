package iterator_test

import (
	"testing"

	"github.com/ndabAP/assocentity/v9/internal/iterator"
)

var testElems = []int{1, 2, 3, 3, 1, 5, 6}

func newTestIterator() *iterator.Iterator[int] {
	return iterator.New(testElems)
}

func TestNavigation(t *testing.T) {
	it := newTestIterator()

	it.Next()
	if it.CurrElem() != testElems[0] {
		t.Errorf("CurrElem() got = %v, want = %v", it.CurrElem(), testElems[0])
	}

	it.Prev()
	if it.CurrElem() != testElems[0] {
		t.Errorf("CurrElem() got = %v, want = %v", it.CurrElem(), testElems[0])
	}

	it.Foward(1)
	if it.CurrElem() != testElems[1] {
		t.Errorf("CurrElem() got = %v, want = %v", it.CurrElem(), testElems[1])
	}

	it.Rewind(1)
	if it.CurrElem() != testElems[0] {
		t.Errorf("CurrElem() got = %v, want = %v", it.CurrElem(), testElems[0])
	}

	it.Reset()
	// We need an independet counter
	i := 0
	for it.Next() {
		if testElems[i] != it.CurrElem() {
			t.Errorf("CurrElem() got = %v, want = %v", it.CurrElem(), testElems[i])
		}
		i++
	}

	it.SetPos(len(testElems))
	i = len(testElems) - 1
	for it.Prev() {
		if testElems[i] != it.CurrElem() {
			t.Errorf("CurrElem() got = %v, want = %v", it.CurrElem(), testElems[i])
		}
		i--
	}
}

func TestInterface(t *testing.T) {
	it := newTestIterator()

	it.SetPos(1)
	if it.CurrElem() != testElems[1] {
		t.Errorf("SetPos(1) got = %v, want = %v", it.CurrElem(), testElems[1])
	}

	it.Reset()
	it.Next()
	if it.CurrElem() != testElems[0] {
		t.Errorf("Reset() got = %v, want = %v", it.CurrElem(), testElems[1])
	}
}
