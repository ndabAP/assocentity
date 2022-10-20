package nlp

import (
	"github.com/ndabAP/assocentity/v8/internal/iterator"
	"github.com/ndabAP/assocentity/v8/tokenize"
)

// NLPPoSDetermer represents the default part of speech determinator
type NLPPoSDetermer struct{ poS tokenize.PoS }

// NewNLPPoSDetermer returns a new default part of speech determinator
func NewNLPPoSDetermer(poS tokenize.PoS) NLPPoSDetermer { return NLPPoSDetermer{poS} }

// DetermPoS deterimantes if a part of speech tag should be deleted. Ignores
// entities
func (dps NLPPoSDetermer) DetermPoS(textTokens []tokenize.Token, entityTokens [][]tokenize.Token) ([]tokenize.Token, error) {
	// If any part of speech, no need to filter
	if dps.poS == tokenize.ANY {
		return textTokens, nil
	}

	var determTokens []tokenize.Token
	// Prepare for generic iterator
	textElems := make(iterator.Elements, len(textTokens))
	for i, v := range textTokens {
		textElems[i] = v
	}

	textIter := iterator.New(textElems)
	for textIter.Next() {
		textIdx := textIter.CurrPos()

		var (
			entityIter  *iterator.Iterator
			isEntity    bool
			nextTextIdx int = textIdx
		)
		for entityIdx := range entityTokens {
			// Prepare for generic iterator
			entityElems := make(iterator.Elements, len(entityTokens[entityIdx]))
			for i, v := range entityTokens[entityIdx] {
				entityElems[i] = v
			}

			entityIter = iterator.New(entityElems)
			// For every entity token
			for entityIter.Next() {
				isEntity = textIter.CurrElem().(tokenize.Token) == entityIter.CurrElem().(tokenize.Token)
				// Check if first text token matches the entity token
				if !isEntity {
					break
				}

				// Check for next text token
				textIter.Next()
			}

			if isEntity {
				entityIter.Reset()
				for entityIter.Next() {
					determTokens = append(determTokens, entityIter.CurrElem().(tokenize.Token))
				}

				// Skip about the tokenized entity length
				nextTextIdx += entityIter.Len() - 1
				// Entity can't occur twice
				break
			}
		}

		textIter.SetPos(nextTextIdx)

		// Entity already added
		if isEntity {
			continue
		}

		if textIter.CurrElem().(tokenize.Token).PoS&dps.poS != 0 {
			determTokens = append(determTokens, textIter.CurrElem().(tokenize.Token))
		}
	}

	return determTokens, nil
}
