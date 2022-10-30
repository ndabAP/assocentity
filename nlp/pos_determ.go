package nlp

import (
	"github.com/ndabAP/assocentity/v9/internal/comp"
	"github.com/ndabAP/assocentity/v9/internal/iterator"
	"github.com/ndabAP/assocentity/v9/tokenize"
)

// NLPPoSDetermer represents the default part of speech determinator
type NLPPoSDetermer struct{ poS tokenize.PoS }

// NewNLPPoSDetermer returns a new default part of speech determinator
func NewNLPPoSDetermer(poS tokenize.PoS) NLPPoSDetermer { return NLPPoSDetermer{poS} }

// DetermPoS deterimantes if a part of speech tag should be kept. It always
// appends entities
func (dps NLPPoSDetermer) DetermPoS(textTokens []tokenize.Token, entityTokens [][]tokenize.Token) []tokenize.Token {
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
		}

		// Non-entity tokens
		if textIter.CurrElem().PoS&dps.poS != 0 {
			determTokens = append(determTokens, textIter.CurrElem())
		}
	}

	return determTokens
}
