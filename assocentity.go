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
	determTokens := psd.DetermPoS(textTokens, entityTokens)

	determTokensIter := iterator.New(determTokens)
	entityTokensIter := iterator.New(entityTokens)

	// Iterate through part of speech determinated text tokens
	for determTokensIter.Next() {

		// Skip about entity positions, if entity

		currDetermTokensPos := determTokensIter.CurrPos()

		var isEntity bool

		for entityTokensIter.Next() {
			entityTokenIter := iterator.New(entityTokensIter.CurrElem())

			// Compare every entity token with part of speech determianted token
			for entityTokenIter.Next() {
				if entityTokenIter.CurrElem() != determTokensIter.CurrElem() {
					isEntity = false
					break
				}

				// Compare with next text token
				determTokensIter.Next()
			}

			if isEntity {
				// If entity, skip about entity positions and cancel loop
				determTokensIter.SetPos(currDetermTokensPos + entityTokenIter.Len())
				goto IS_ENTITY
			}
		}

		// Reset state if no entity
		determTokensIter.SetPos(currDetermTokensPos)

	IS_ENTITY:

		// Distance
		var entityDist float64

		// Iterate in positive and negative direction to find entity distances
		posDirIter := determTokensIter
		negDirIter := determTokensIter

		for posDirIter.Next() {
			currPosDirPos := negDirIter.CurrPos()

			// Tells if current text token is entity and how many positions from
			// here
			isEntity, entityLen := entityChecker(posDirIter, entityTokens, posDir)
			if isEntity {
				appendMap(assocEntitiesVals, determTokensIter.CurrElem().Token, entityDist)
				// Skip about entity
				posDirIter.SetPos(currPosDirPos + entityLen)
			}

			entityDist++

			if isEntity {
				continue
			}

			// Reset because mutated
			posDirIter.SetPos(currPosDirPos)
		}

		// Reset distance
		entityDist = 0

		for negDirIter.Prev() {
			negDirIdx := negDirIter.CurrPos()

			isEntity, entityLen := entityChecker(negDirIter, entityTokens, negDir)
			if isEntity {
				appendMap(assocEntitiesVals, determTokensIter.CurrElem().Token, entityDist)

				negDirIter.SetPos(negDirIdx - entityLen)
			}

			entityDist++

			if isEntity {
				continue
			}

			negDirIter.SetPos(negDirIdx)
		}
	}

	// Calculate the average distances
	for token, dist := range assocEntitiesVals {
		assocEntities[token] = avg(dist)
	}
	return assocEntities, err
}

// Iterates through entity and PoS determinated tokens and returns true if found
// and positions to skip
func entityChecker(determTokensIter *iterator.Iterator[tokenize.Token], entityTokens [][]tokenize.Token, dir direction) (bool, int) {
	var isEntity bool
	for entityIdx := range entityTokens {
		entityIter := iterator.New(entityTokens[entityIdx])

		switch dir {

		case posDir:
			for entityIter.Next() {
				if determTokensIter.CurrElem() != entityIter.CurrElem() {
					// Check if first token matches the entity token
					isEntity = false
					break
				}

				// Check for next token
				determTokensIter.Next()
			}

		case negDir:
			for entityIter.Prev() {
				if determTokensIter.CurrElem() != entityIter.CurrElem() {
					// Check if first token matches the entity token
					isEntity = false
					break
				}

				// Check for next token
				determTokensIter.Prev()
			}
		}

		if isEntity {
			return isEntity, entityIter.Len()
		}
	}

	return false, 0
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
