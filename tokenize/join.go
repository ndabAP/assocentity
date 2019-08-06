package tokenize

import (
	"strings"

	"github.com/ndabAP/assocentity/v3/internal/iterator"
)

// Joiner joines a tokenizer
type Joiner interface {
	Join(Tokenizer) error
	Tokens() []string
}

// DefaultJoin is the default joiner
type DefaultJoin struct {
	tokens []string
	sep    string // Separator
}

// NewDefaultJoin returns a new default join
func NewDefaultJoin(sep string) *DefaultJoin {
	return &DefaultJoin{[]string{}, sep}
}

// Join joins strings in a string slice
func (dj *DefaultJoin) Join(tok Tokenizer) error {
	textTokens, err := tok.TokenizeText()
	if err != nil {
		return err
	}

	entityTokens, err := tok.TokenizeEntities()
	if err != nil {
		return err
	}

	textTraverser := iterator.New(&textTokens)
	// For every text token
	for textTraverser.Next() {
		textIdx := textTraverser.CurrPos()
		// For every tokenized entity
		for entityIdx := range entityTokens {
			entityTraverser := iterator.New(&entityTokens[entityIdx])

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

				// Check for next text token
				textTraverser.Next()
			}

			if isEntity {
				// Merge the entity
				textTokens[textIdx] = strings.Join(entityTokens[entityIdx], dj.sep)
				// Remove text tokens that contain the entity
				idx := textIdx + 1
				textTokens = append(textTokens[:idx], textTokens[idx+entityTraverser.Len()-1:]...)
			}
		}

		textTraverser.SetPos(textIdx)
	}

	dj.tokens = textTokens

	return nil
}

// Tokens returns the joined tokens
func (dj *DefaultJoin) Tokens() []string {
	return dj.tokens
}
