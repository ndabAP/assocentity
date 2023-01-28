package pos

import (
	"github.com/ndabAP/assocentity/v12/internal/comp"
	"github.com/ndabAP/assocentity/v12/internal/iterator"
	"github.com/ndabAP/assocentity/v12/tokenize"
)

// poSDetermer represents the default part of speech determinator
type poSDetermer struct{ poS tokenize.PoS }

// NewPoSDetermer returns a new default part of speech determinator
func NewPoSDetermer(poS tokenize.PoS) poSDetermer { return poSDetermer{poS} }

// DetermPoS deterimantes if a part of speech tag should be kept. It always
// appends entities
func (dps poSDetermer) DetermPoS(textTokens []tokenize.Token, entityTokens [][]tokenize.Token) []tokenize.Token {
	// If any part of speech, no need to determinate
	if dps.poS == tokenize.ANY {
		return textTokens
	}

	var determTokens []tokenize.Token

	textIter := iterator.New(textTokens)
	entityTokensIter := iterator.New(entityTokens)

	for textIter.Next() {
		currTextPos := textIter.CurrPos()
		isEntity, entity := comp.TextWithEntities(textIter, entityTokensIter, comp.DirPos)
		if isEntity {
			textIter.SetPos(currTextPos + len(entity))
			// Entity is always kept
			determTokens = append(determTokens, entity...)
			continue
		}

		// Non-entity tokens
		if textIter.CurrElem().PoS&dps.poS != 0 {
			determTokens = append(determTokens, textIter.CurrElem())
		}
	}

	return determTokens
}
