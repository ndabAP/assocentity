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
func TextWithEntity(textIter *iterator.Iterator[tokenize.Token], entityTokensIter *iterator.Iterator[[]tokenize.Token], entityIterDir Direction) (bool, []tokenize.Token) {
	// Reset iterators position after comparing
	currTextPos := textIter.CurrPos()
	defer textIter.SetPos(currTextPos)
	currEntityTokensPos := entityTokensIter.CurrPos()
	defer entityTokensIter.SetPos(currEntityTokensPos)

	var isEntity bool = true
	for entityTokensIter.Next() {
		entityIter := iterator.New(entityTokensIter.CurrElem())

		switch entityIterDir {

		// ->
		case DirPos:
			for entityIter.Next() {
				if textIter.CurrElem() != entityIter.CurrElem() {
					// Check if first token matches the entity token
					isEntity = false
					break
				}

				// Check for next token
				textIter.Next()
			}

		// <-
		case DirNeg:
			for entityIter.Prev() {
				if textIter.CurrElem() != entityIter.CurrElem() {
					isEntity = false
					break
				}

				textIter.Prev()
			}
		}

		if isEntity {
			return true, entityTokensIter.CurrElem()
		}
	}

	return false, []tokenize.Token{}
}
