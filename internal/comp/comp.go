package comp

import (
	"github.com/ndabAP/assocentity/v11/internal/iterator"
	"github.com/ndabAP/assocentity/v11/tokenize"
)

type Direction int

var (
	DirPos Direction = 1
	DirNeg Direction = -1
)

// Checks if current text token is entity and if, returns entity
func TextWithEntities(textIter *iterator.Iterator[tokenize.Token], entityTokensIter *iterator.Iterator[[]tokenize.Token], entityIterDir Direction) (bool, []tokenize.Token) {
	// Reset iterators position after comparing (and before)
	currTextPos := textIter.CurrPos()
	defer entityTokensIter.Reset()
	defer textIter.SetPos(currTextPos)
	entityTokensIter.Reset()

	// By default, we assume an entity
	var isEntity bool = true

	for entityTokensIter.Next() {
		// Reset
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

				// Advance text iterator to compare against
				textIter.Next()
			}

		// <-
		case DirNeg:
			// We scan backwards and start from top
			entityIter.SetPos(entityIter.Len()) // [1, 2, 3, (4)]
			for entityIter.Prev() {
				if !eqItersElems(textIter, entityIter) {
					isEntity = false
				}

				textIter.Prev()
			}
		}

		if isEntity {
			return true, entityTokensIter.CurrElem()
		} else {
			// Reset to compare with next entity tokens
			textIter.SetPos(currTextPos)
		}
	}

	return false, []tokenize.Token{}
}

func eqItersElems(x *iterator.Iterator[tokenize.Token], y *iterator.Iterator[tokenize.Token]) bool {
	return x.CurrElem() == y.CurrElem()
}
