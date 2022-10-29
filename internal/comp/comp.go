package comp

import (
	"github.com/ndabAP/assocentity/v8/internal/iterator"
	"github.com/ndabAP/assocentity/v8/tokenize"
)

type Direction int

var (
	DirPos Direction = 1
	DirNeg Direction = -1
)

// Checks if current text token is entity token
func TextWithEntities(textIter *iterator.Iterator[tokenize.Token], entityTokensIter *iterator.Iterator[[]tokenize.Token], entityIterDir Direction) (bool, []tokenize.Token) {
	// Reset iterators position after comparing
	currTextPos := textIter.CurrPos()

	defer entityTokensIter.Reset()
	entityTokensIter.Reset()

	var isEntity bool = true

	for entityTokensIter.Next() {
		isEntity = true

		entityIter := iterator.New(entityTokensIter.CurrElem())

		switch entityIterDir {

		// ->
		case DirPos:
			for entityIter.Next() {
				// Check if text token matches the entity token
				if !eqItersElems(textIter, entityIter) {
					isEntity = false
				}

				textIter.Next()
			}

		// <-
		case DirNeg:
			// We scan backwards
			entityIter.SetPos(entityIter.Len()) // [1, 2, 3, (4)]
			for entityIter.Prev() {
				if !eqItersElems(textIter, entityIter) {
					isEntity = false
				}

				textIter.Prev()
			}
		}

		textIter.SetPos(currTextPos)

		if isEntity {
			return true, entityTokensIter.CurrElem()
		}
	}

	textIter.SetPos(currTextPos)

	return false, []tokenize.Token{}
}

func eqItersElems(x *iterator.Iterator[tokenize.Token], y *iterator.Iterator[tokenize.Token]) bool {
	return x.CurrElem() == y.CurrElem()
}
