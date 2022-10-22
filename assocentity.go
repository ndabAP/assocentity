package assocentity

import (
	"context"

	"github.com/ndabAP/assocentity/v8/internal/iterator"
	"github.com/ndabAP/assocentity/v8/tokenize"
)

type direction int

var (
	posDir direction = 1
	negDir direction = -1
)

func Do(ctx context.Context, tokenizer tokenize.Tokenizer, psd tokenize.PoSDetermer, text string, entities []string) (map[string]float64, error) {
	var (
		assocEntities     = make(map[string]float64)
		assocEntitiesVals = make(map[string][]float64)

		err error
	)

	// Tokenize text
	textTokens, err := tokenizer.Tokenize(ctx, text)
	if err != nil {
		return assocEntities, err
	}

	// Tokenize entites
	var entityTokens [][]tokenize.Token
	for _, entity := range entities {
		tokens, err := tokenizer.Tokenize(ctx, entity)
		if err != nil {
			return assocEntities, err
		}
		entityTokens = append(entityTokens, tokens)
	}

	// Determinate part of speech
	determTokens, err := psd.DetermPoS(textTokens, entityTokens)
	if err != nil {
		return assocEntities, err
	}

	// Create iterators
	determTokensIter := iterator.New[tokenize.Token](determTokens)
	entityTokensIter := iterator.New[[]tokenize.Token](entityTokens)

	// Iterate through part of speech determinated text tokens
	for determTokensIter.Next() {
		// Skip about entity positions if entity. It's easier to do that up
		// front instead of during positive/negative iteration

		currDetermTokensPos := determTokensIter.CurrPos()

		var isEntity bool

		// For every entity tokens
		for entityTokensIter.Next() {
			entityTokenIter := iterator.New[tokenize.Token](entityTokensIter.CurrElem())

			// Compare every entity token with part of speech determianted token
			for entityTokenIter.Next() {
				if entityTokenIter.CurrElem() != determTokensIter.CurrElem() {
					isEntity = false
				}

				determTokensIter.Next()
			}

			if isEntity {
				// If entity, skip about entity positions
				determTokensIter.SetPos(currDetermTokensPos + entityTokenIter.Len())
				goto SKIP_ENTITY
			}
		}

		// Reset state if no entity
		determTokensIter.SetPos(currDetermTokensPos)

	SKIP_ENTITY:

		// Distance
		var entityDistance float64

		// Iterate in positive and negative direction to find entity distances
		posDirIter := determTokensIter
		negDirIter := determTokensIter

		for posDirIter.Next() {
			posDirIdx := posDirIter.CurrPos()

			isEntity, pos := entityChecker(posDirIter, entityTokens, posDir)
			if isEntity {
				appendMap(assocEntitiesVals, determTokensIter.CurrElem().Token, entityDistance)
				// Skip about entity
				posDirIter.SetPos(posDirIdx + pos)
			}

			entityDistance++

			if isEntity {
				continue
			}

			// Reset because mutated
			posDirIter.SetPos(posDirIdx)
		}

		// Reset distance
		entityDistance = 0

		for negDirIter.Prev() {
			negDirIdx := negDirIter.CurrPos()

			isEntity, pos := entityChecker(negDirIter, entityTokens, negDir)
			if isEntity {
				appendMap(assocEntitiesVals, determTokensIter.CurrElem().Token, entityDistance)

				negDirIter.SetPos(negDirIdx - pos)
			}

			entityDistance++

			if isEntity {
				continue
			}

			negDirIter.SetPos(negDirIdx)
		}
	}

	// Calculate the distances
	for elem, dist := range assocEntitiesVals {
		assocEntities[elem] = avg(dist)
	}
	return assocEntities, err
}

// Iterates through entity and PoS determinated tokens and returns true if found
// and positions to skip
func entityChecker[T any](determTokensIter *iterator.Iterator[tokenize.Token], entityTokens [][]tokenize.Token, dir direction) (bool, int) {
	var (
		entityIter *iterator.Iterator
		isEntity   bool
	)
	for entityIdx := range entityTokens {
		// Prepare for generic iterator
		entityElems := make(iterator.Elements, len(entityTokens[entityIdx]))
		for i, v := range entityTokens[entityIdx] {
			entityElems[i] = v
		}

		entityIter = iterator.New(entityElems)
		switch dir {

		case posDir:
			// Positive direction
			for entityIter.Next() {
				isEntity = determTokensIter.CurrElem().(tokenize.Token).Token == entityIter.CurrElem().(tokenize.Token).Token
				// Check if first token matches the entity token
				if !isEntity {
					break
				}

				// Check for next token
				determTokensIter.Next()
			}

		case negDir:
			// Negative direction

			entityIter.SetPos(entityIter.Len() - 1)

			for entityIter.Prev() {
				isEntity = determTokensIter.CurrElem().(tokenize.Token).Token == entityIter.CurrElem().(tokenize.Token).Token
				// Check if first token matches the entity token
				if !isEntity {
					break
				}

				// Check for next token
				determTokensIter.Prev()
			}
		}

		if isEntity {
			return isEntity, entityIter.Len() - 1
		}
	}

	return isEntity, entityIter.Len() - 1
}

// Returns the average of a float slice
func avg(xs []float64) float64 {
	total := 0.0
	for _, v := range xs {
		total += v
	}
	return total / float64(len(xs))
}

// Helper to append float to a map
func appendMap(m map[string][]float64, k string, f float64) {
	m[k] = append(m[k], f)
}
