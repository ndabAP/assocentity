package tokenize

import "github.com/ndabAP/assocentity/v3/internal/generator"

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
		entityTraverser := generator.New(entityTokens[textTraverser.GetCurrPos()])

		var isEntity bool
		// For every entity token
		for entityTraverser.Next() {
			// Check if first text token matches the entity token
			if textTraverser.GetCurrElem() != entityTraverser.GetCurrElem() {
				break
			}

		}
	}

	return []string{}, nil
}
