package tokenize

import (
	"strings"

	"github.com/ndabAP/assocentity/v3/internal/generator"
)

// Joiner joines a tokenizer
type Joiner interface {
	Join(Tokenizer) ([]string, error)
}

// DefaultJoin is the default joiner
type DefaultJoin struct {
	sep string // Seperator
}

// NewDefaultJoin returns a new default join
func NewDefaultJoin(sep string) DefaultJoin {
	return DefaultJoin{sep}
}

// Join joins strings in a string slice
func (dm *DefaultJoin) Join(tok Tokenizer) ([]string, error) {
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
				idx := currTextIdx + 1
				textTokens = append(textTokens[:idx], textTokens[idx+entityTraverser.Len()-1:]...)
			}
		}

		textTraverser.SetPos(currTextIdx)
	}

	return textTokens, nil
}
