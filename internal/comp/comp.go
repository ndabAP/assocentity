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
	defer textIter.SetPos(currTextPos)
	defer entityTokensIter.Reset()

	entityTokensIter.Reset()

	var isEntity bool = true
	for entityTokensIter.Next() {
		entityIter := iterator.New(entityTokensIter.CurrElem())

		switch entityIterDir {

		// ->
		case DirPos:
			for entityIter.Next() {
				// Check if text token matches the entity token and advance
				// one token
				if (textIter.CurrElem() != entityIter.CurrElem()) || !textIter.Next() {
					isEntity = false
					break
				}
			}

		// <-
		case DirNeg:
			// We scan backwards
			entityIter.SetPos(entityIter.Len()) // [1, 2, 3, (4)]
			for entityIter.Prev() {
				if (textIter.CurrElem() != entityIter.CurrElem()) || !textIter.Prev() {
					isEntity = false
					break
				}
			}
		}

		if isEntity {
			return true, entityTokensIter.CurrElem()
		}
	}

	return false, []tokenize.Token{}
}

func Iters(x *iterator.Iterator[tokenize.Token], y *iterator.Iterator[tokenize.Token]) bool {
	return x.CurrElem() == y.CurrElem()
}
