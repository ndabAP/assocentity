package assocentity

import (
	"context"
	"math"

	"github.com/ndabAP/assocentity/v8/internal/iterator"
	"github.com/ndabAP/assocentity/v8/tokenize"
)

type direction int

var (
	posDir direction = 1
	negDir direction = -1
)

// Do returns the entity distances
func Do(ctx context.Context, tokenizer tokenize.Tokenizer, psd tokenize.PoSDetermer, text string, entities []string) (map[string]float64, error) {
	var (
		assocEntities = make(map[string]float64)
		distAccum     = make(map[string][]float64)

		err error
	)

	textTokens, err := tokenizer.Tokenize(ctx, text)
	if err != nil {
		return assocEntities, err
	}

	var entityTokens [][]tokenize.Token
	for _, entity := range entities {
		tokenized, err := tokenizer.Tokenize(ctx, entity)
		if err != nil {
			return assocEntities, err
		}
		entityTokens = append(entityTokens, tokenized)
	}

	determTokens, err := psd.DetermPoS(textTokens, entityTokens)
	if err != nil {
		return assocEntities, err
	}

	// Prepare for generic iterator
	determElems := make(iterator.Elements, len(determTokens))
	for i, v := range determTokens {
		determElems[i] = v
	}

	determTokensIter := iterator.New(determElems)
	for determTokensIter.Next() {
		determTokensIdx := determTokensIter.CurrPos()

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
			for entityIter.Next() {
				isEntity = determTokensIter.CurrElem().(tokenize.Token) == entityIter.CurrElem().(tokenize.Token)
				// Check if first token matches the entity token
				if !isEntity {
					break
				}

				// Check for next token
				determTokensIter.Next()
			}

			if isEntity {
				break
			}
		}

		if isEntity {
			// Skip about entity positions
			determTokensIter.SetPos(determTokensIdx + entityIter.Len() - 1)
			continue
		}

		// Reset because mutated
		determTokensIter.SetPos(determTokensIdx)

		// Iterate in positive and negative direction

		// Distance
		var entityDistance float64

		// Iterate positive direction
		posDirIter := iterator.New(determElems)
		posDirIter.SetPos(determTokensIdx)
		for posDirIter.Next() {
			posDirIdx := posDirIter.CurrPos()

			isEntity, len := entityChecker(posDirIter, entityTokens, posDir)
			if isEntity {
				distAccum[determTokensIter.CurrElem().(tokenize.Token).Token] = append(distAccum[determTokensIter.CurrElem().(tokenize.Token).Token], entityDistance)
				// Skip about entity
				posDirIter.SetPos(posDirIdx + len - 1)
			}

			entityDistance++

			if isEntity {
				continue
			}

			// Reset because entityChecker is mutating
			posDirIter.SetPos(posDirIdx)
		}

		// Iterate negative direction
		// Reset distance
		entityDistance = 0

		negDirIter := iterator.New(determElems)
		negDirIter.SetPos(determTokensIdx)
		for negDirIter.Prev() {
			negDirIdx := negDirIter.CurrPos()

			isEntity, len := entityChecker(negDirIter, entityTokens, negDir)
			if isEntity {
				distAccum[determTokensIter.CurrElem().(tokenize.Token).Token] = append(distAccum[determTokensIter.CurrElem().(tokenize.Token).Token], entityDistance)
				// Skip about entity
				negDirIter.SetPos(negDirIdx - len + 1)
			}

			entityDistance++

			if isEntity {
				continue
			}

			// Reset because entityChecker is mutating
			negDirIter.SetPos(negDirIdx)
		}
	}

	// Calculate the distances
	for elem, dist := range distAccum {
		assocEntities[elem] = avg(dist)
	}
	return assocEntities, nil
}

// Iterates through entity and PoS determinated tokens and returns true if found
// and positions to skip
func entityChecker(determTokTraverser *iterator.Iterator, entityTokens [][]tokenize.Token, dir direction) (bool, int) {
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
				isEntity = determTokTraverser.CurrElem().(tokenize.Token).Token == entityIter.CurrElem().(tokenize.Token).Token
				// Check if first token matches the entity token
				if !isEntity {
					break
				}

				// Check for next token
				determTokTraverser.Next()
			}

		case negDir:
			// Negative direction

			entityIter.SetPos(entityIter.Len() - 1)

			for entityIter.Prev() {
				isEntity = determTokTraverser.CurrElem().(tokenize.Token).Token == entityIter.CurrElem().(tokenize.Token).Token
				// Check if first token matches the entity token
				if !isEntity {
					break
				}

				// Check for next token
				determTokTraverser.Prev()
			}
		}

		if isEntity {
			return isEntity, entityIter.Len()
		}
	}

	return isEntity, entityIter.Len()
}

// Returns the average of a float slice
func avg(xs []float64) float64 {
	total := 0.0
	for _, v := range xs {
		total += v
	}
	return round(total / float64(len(xs)))
}

// Rounds to nearest 0.5
func round(x float64) float64 {
	return math.Round(x/0.5) * 0.5
}
