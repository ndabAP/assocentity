package tokenize

import (
	"strings"

	"github.com/ndabAP/assocentity/v5/internal/iterator"
)

// Joiner joines a tokenizer taking the part of speech determinator into account
type Joiner interface {
	Join(PoSDetermer, Tokenizer) ([]string, error)
}

// Join is the default joiner
type Join struct {
	sep string // Separator
}

// NewJoin returns a new default join
func NewJoin(sep string) *Join { return &Join{sep} }

const (
	// Whitespace is the default separator
	Whitespace = " "
)

// Join joins strings in a string slice
func (dj *Join) Join(dps PoSDetermer, tokenizer Tokenizer) ([]string, error) {
	textTokens, err := dps.Determ(tokenizer)
	if err != nil {
		return []string{}, err
	}

	entityTokens, err := tokenizer.TokenizeEntities()
	if err != nil {
		return []string{}, err
	}

	var res []string
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
			// Skip single value entities
			if entityTraverser.Len() == 1 {
				break
			}

			for entityTraverser.Next() {
				isEntity = textTraverser.CurrElem().(string) == entityTraverser.CurrElem().(Token).Token
				// Check if first text token matches the entity token
				if !isEntity {
					break
				}

				// Check for next text token
				textTraverser.Next()
			}

			if isEntity {
				entityTraverser.Reset()
				var entity []string
				for entityTraverser.Next() {
					entity = append(entity, entityTraverser.CurrElem().(Token).Token)
				}

				// Merge the entity
				res = append(res, strings.Join(entity, dj.sep))

				// Skip about the tokenized entity length
				nextTextTraverserPos += entityTraverser.Len() - 1
				// Entity can't occur twice
				break
			}
		}

		textTraverser.SetPos(nextTextTraverserPos)

		// Entity is already joined
		if isEntity {
			continue
		}

		res = append(res, textTraverser.CurrElem().(string))
	}

	return res, nil
}
