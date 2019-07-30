package tokenize

import (
	"strings"

	"github.com/ndabAP/assocentity/v3/internal/generator"
)

// Multiplexer multiplexes a tokenizer
type Multiplexer interface {
	Multiplex(Tokenizer) ([]string, error)
}

// DefaultMultiplex is the default multiplexer
type DefaultMultiplex struct{}

// NewDefaultMultiplex returns a new default multiplex
func NewDefaultMultiplex() DefaultMultiplex {
	return DefaultMultiplex{}
}

// Multiplex multiplexes a string slice
func (dm *DefaultMultiplex) Multiplex(tok Tokenizer) ([]string, error) {
	textTokens, err := tok.TokenizeText()
	if err != nil {
		return nil, err
	}

	entityTokens, err := tok.TokenizeEntities()
	if err != nil {
		return nil, err
	}

	textTraverser := generator.New(textTokens)
	// For every text token
	for textTraverser.Next() {
		currTextIdx := textTraverser.CurrPos()
		// For every tokenized entity
		for idx := range entityTokens {
			currEntityIdx := idx
			entityTraverser := generator.New(entityTokens[currEntityIdx])

			// Skip single value entities
			if entityTraverser.Len() == 1 {
				break
			}

			var isEntity bool
			// For every entity token
			for entityTraverser.Next() {
				isEntity = textTraverser.CurrElem() == entityTraverser.CurrElem()
				// Check if first text token matches the entity token
				if !isEntity {
					break
				}

				isEntity = textTraverser.CurrElem() == entityTraverser.CurrElem()
				// Check for next text token
				textTraverser.Next()

				if !isEntity {
					break
				}
			}

			if isEntity {
				// Merge the entity
				textTokens[currTextIdx] = strings.Join(entityTokens[currEntityIdx], " ")
				// Remove text tokens that contain the entity
				entityTraverser.Reset()
				for entityTraverser.Next() {
					idx := currTextIdx + 1
					textTokens = append(textTokens[:idx], textTokens[idx+1:]...)
				}
			}
		}

		textTraverser.SetPos(currTextIdx)
	}

	return textTokens, nil
}
