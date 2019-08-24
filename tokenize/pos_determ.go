package tokenize

import (
	"github.com/ndabAP/assocentity/v5/internal/iterator"
)

// PoSDetermer determinates if part of speech tags should be deleted
type PoSDetermer interface {
	Determ(Tokenizer) ([]Token, error)
}

// PoSDeterm represents the default part of speech determinator
type PoSDeterm struct{ poS int }

// NewPoSDetermer returns a new default part of speech determinator
func NewPoSDetermer(poS int) *PoSDeterm { return &PoSDeterm{poS} }

// Determ deterimantes if a part of speech tag should be deleted
func (dps *PoSDeterm) Determ(tokenizer Tokenizer) ([]Token, error) {
	textTokens, err := tokenizer.TokenizeText()
	if err != nil {
		return []Token{}, err
	}

	entityTokens, err := tokenizer.TokenizeEntities()
	if err != nil {
		return []Token{}, err
	}

	var res []Token
	// Prepare for generic iterator
	t := make(iterator.Elements, len(textTokens))
	for i, v := range textTokens {
		t[i] = v
	}

	textTraverser := iterator.New(t)
	for textTraverser.Next() {
		textIdx := textTraverser.CurrPos()

		var (
			isEntity             bool
			entityTraverser      *iterator.Iterator
			nextTextTraverserPos int = textIdx
		)
		for entityIdx := range entityTokens {
			// Prepare for generic iterator
			e := make(iterator.Elements, len(entityTokens[entityIdx]))
			for i, v := range entityTokens[entityIdx] {
				e[i] = v
			}

			entityTraverser = iterator.New(e)
			// For every entity token
			for entityTraverser.Next() {
				isEntity = textTraverser.CurrElem().(Token) == entityTraverser.CurrElem().(Token)
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
					res = append(res, entityTraverser.CurrElem().(Token))
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

		if textTraverser.CurrElem().(Token).PoS&dps.poS != 0 {
			res = append(res, textTraverser.CurrElem().(Token))
		}
	}

	return res, nil
}
