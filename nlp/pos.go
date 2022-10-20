package nlp

import (
	"github.com/ndabAP/assocentity/v8/internal/iterator"
	"github.com/ndabAP/assocentity/v8/tokenize"
)

// PoSDeterm represents the default part of speech determinator
type PoSDeterm struct{ poS tokenize.PoS }

// NewPoSDetermer returns a new default part of speech determinator
func NewPoSDetermer(poS tokenize.PoS) *PoSDeterm { return &PoSDeterm{poS} }

// Determ deterimantes if a part of speech tag should be deleted
func (dps *PoSDeterm) Determ(tokenizedText []tokenize.Token, tokenizedEntities [][]tokenize.Token) ([]tokenize.Token, error) {
	// If any part of speech, no need to filter
	if dps.poS == tokenize.ANY {
		return tokenizedText, nil
	}

	var tokens []tokenize.Token
	// Prepare for generic iterator
	t := make(iterator.Elements, len(tokenizedText))
	for i, v := range tokenizedText {
		t[i] = v
	}

	textTraverser := iterator.New(t)
	for textTraverser.Next() {
		textIdx := textTraverser.CurrPos()

		var (
			entityTraverser      *iterator.Iterator
			isEntity             bool
			nextTextTraverserPos int = textIdx
		)
		for entityIdx := range tokenizedEntities {
			// Prepare for generic iterator
			elems := make(iterator.Elements, len(tokenizedEntities[entityIdx]))
			for i, v := range tokenizedEntities[entityIdx] {
				elems[i] = v
			}

			entityTraverser = iterator.New(elems)
			// For every entity token
			for entityTraverser.Next() {
				isEntity = textTraverser.CurrElem().(tokenize.Token) == entityTraverser.CurrElem().(tokenize.Token)
				// Check if first text token matches the entity token
				if !isEntity {
					break
				}

				// Check for next text token
				textTraverser.Next()
			}

			if isEntity {
				entityTraverser.Reset()
				for entityTraverser.Next() {
					tokens = append(tokens, entityTraverser.CurrElem().(tokenize.Token))
				}

				// Skip about the tokenized entity length
				nextTextTraverserPos += entityTraverser.Len() - 1
				// Entity can't occur twice
				break
			}
		}

		textTraverser.SetPos(nextTextTraverserPos)

		// Entity already added
		if isEntity {
			continue
		}

		if textTraverser.CurrElem().(tokenize.Token).PoS&dps.poS != 0 {
			tokens = append(tokens, textTraverser.CurrElem().(tokenize.Token))
		}
	}

	return tokens, nil
}
