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
func Do(ctx context.Context, tokenizer tokenize.Tokenizer, psd tokenize.PoSDetermer, text string, entities []string) (map[tokenize.Token]float64, error) {
	textTokens, err := tokenizer.Tokenize(ctx, text)
	if err != nil {
		return map[tokenize.Token]float64{}, err
	}

	var entityTokens [][]tokenize.Token
	for _, entity := range entities {
		tokenized, err := tokenizer.Tokenize(ctx, entity)
		if err != nil {
			return map[tokenize.Token]float64{}, err
		}

		entityTokens = append(entityTokens, tokenized)
	}

	determTokens, err := psd.Determ(textTokens, entityTokens)
	if err != nil {
		return map[tokenize.Token]float64{}, err
	}

	// Prepare for generic iterator
	determElems := make(iterator.Elements, len(determTokens))
	for i, v := range determTokens {
		determElems[i] = v
	}

	var distAccum = make(map[tokenize.Token][]float64)

	determTokIter := iterator.New(determElems)
	for determTokIter.Next() {
		determTokIdx := determTokIter.CurrPos()

		var (
			entityIter *iterator.Iterator
			isEntity   bool
		)
		for entityIdx := range entityTokens {
			// Prepare for generic iterator
			e := make(iterator.Elements, len(entityTokens[entityIdx]))
			for i, v := range entityTokens[entityIdx] {
				e[i] = v
			}

			entityIter = iterator.New(e)
			for entityIter.Next() {
				isEntity = determTokIter.CurrElem().(tokenize.Token) == entityIter.CurrElem().(tokenize.Token)
				// Check if first token matches the entity token
				if !isEntity {
					break
				}

				// Check for next token
				determTokIter.Next()
			}

			if isEntity {
				break
			}
		}

		if isEntity {
			// Skip about entity positions
			determTokIter.SetPos(determTokIdx + entityIter.Len() - 1)
			continue
		}

		// Reset because mutated
		determTokIter.SetPos(determTokIdx)

		// Distance
		var dist float64

		// Iterate positive direction
		posDirIter := iterator.New(determElems)
		posDirIter.SetPos(determTokIdx)
		for posDirIter.Next() {
			posDirIdx := posDirIter.CurrPos()

			isEntity, len := entityChecker(posDirIter, entityTokens, posDir)
			if isEntity {
				distAccum[determTokIter.CurrElem().(tokenize.Token)] = append(distAccum[determTokIter.CurrElem().(tokenize.Token)], dist)
				// Skip about entity
				posDirIter.SetPos(posDirIdx + len - 1)
			}

			dist++

			if isEntity {
				continue
			}

			// Reset because entityChecker is mutating
			posDirIter.SetPos(posDirIdx)
		}

		// Iterate negative direction
		// Reset distance
		dist = 0
		negDirIter := iterator.New(determElems)
		negDirIter.SetPos(determTokIdx)
		for negDirIter.Prev() {
			negDirIdx := negDirIter.CurrPos()

			isEntity, len := entityChecker(negDirIter, entityTokens, negDir)
			if isEntity {
				distAccum[determTokIter.CurrElem().(tokenize.Token)] = append(distAccum[determTokIter.CurrElem().(tokenize.Token)], dist)
				// Skip about entity
				negDirIter.SetPos(negDirIdx - len + 1)
			}

			dist++

			if isEntity {
				continue
			}

			// Reset because entityChecker is mutating
			negDirIter.SetPos(negDirIdx)
		}
	}

	assocEntities := make(map[tokenize.Token]float64)
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
		elems := make(iterator.Elements, len(entityTokens[entityIdx]))
		for i, v := range entityTokens[entityIdx] {
			elems[i] = v
		}

		entityIter = iterator.New(elems)
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
