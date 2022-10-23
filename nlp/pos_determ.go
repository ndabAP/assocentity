package nlp

import (
	"github.com/ndabAP/assocentity/v8/internal/iterator"
	"github.com/ndabAP/assocentity/v8/tokenize"
)

// NLPPoSDetermer represents the default part of speech determinator
type NLPPoSDetermer struct{ poS tokenize.PoS }

// NewNLPPoSDetermer returns a new default part of speech determinator
func NewNLPPoSDetermer(poS tokenize.PoS) NLPPoSDetermer { return NLPPoSDetermer{poS} }

// DetermPoS deterimantes if a part of speech tag should be deleted. It appends
// entities without filtering
func (dps NLPPoSDetermer) DetermPoS(textTokens []tokenize.Token, entityTokens [][]tokenize.Token) []tokenize.Token {
	// If any part of speech, no need to determinate
	if dps.poS == tokenize.ANY {
		return textTokens
	}

	var determTokens []tokenize.Token

	textIter := iterator.New(textTokens)
	entityTokensIter := iterator.New(entityTokens)

IS_ENTITY:
	for textIter.Next() {
		// Reset from previous iteration
		entityTokensIter.Reset()

		currTextPos := textIter.CurrPos()

		var isEntity bool

		for entityTokensIter.Next() {
			entityTokenIter := iterator.New(entityTokensIter.CurrElem())

			// Compare every entity token with part of speech determianted token
			for entityTokenIter.Next() {
				if textIter.CurrElem() != entityTokenIter.CurrElem() {
					isEntity = false
					break
				}

				// Compare with next text token
				textIter.Next()
			}

			if isEntity {
				// Append entity without filtering
				for entityTokenIter.Next() {
					determTokens = append(determTokens, entityTokenIter.CurrElem())
				}

				// If entity, skip about entity positions and cancel loop
				textIter.SetPos(currTextPos + entityTokenIter.Len())
				goto IS_ENTITY
			}
		}

		// Reset state if no entity
		// TODO!: Set pos sets init to true even its not
		textIter.SetPos(currTextPos)

		// Non-entity tokens
		if textIter.CurrElem().PoS&dps.poS != 0 {
			determTokens = append(determTokens, textIter.CurrElem())
		}
	}

	return determTokens
}
