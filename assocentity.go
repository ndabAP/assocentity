package assocentity

import (
	"context"

	"github.com/ndabAP/assocentity/v8/internal/comp"
	"github.com/ndabAP/assocentity/v8/internal/iterator"
	"github.com/ndabAP/assocentity/v8/tokenize"
)

func Do(ctx context.Context, tokenizer tokenize.Tokenizer, psd tokenize.PoSDetermer, text string, entities []string) (map[string]float64, error) {
	var (
		assocTokens      = make(map[string]float64)
		assocTokensAccum = make(map[string][]float64)

		err error
	)

	// Tokenize text
	textTokens, err := tokenizer.Tokenize(ctx, text)
	if err != nil {
		return assocTokens, err
	}

	// Tokenize entites
	var entityTokens [][]tokenize.Token
	for _, entity := range entities {
		tokens, err := tokenizer.Tokenize(ctx, entity)
		if err != nil {
			return assocTokens, err
		}
		entityTokens = append(entityTokens, tokens)
	}

	// Determinate part of speech
	determTokens := psd.DetermPoS(textTokens, entityTokens)

	determTokensIter := iterator.New(determTokens)
	entityTokensIter := iterator.New(entityTokens)

	// Iterate through part of speech determinated text tokens
	for determTokensIter.Next() {
		// If the current token is an entity, we skip about the entity
		currDetermTokensPos := determTokensIter.CurrPos()
		isEntity, entity := comp.TextWithEntities(determTokensIter, entityTokensIter, comp.DirPos)
		if isEntity {
			determTokensIter.Foward(len(entity) - 1)
			continue
		}

		// Now we can collect the actual distances

		// Distance
		var entityDist float64

		// TODO: Create once and store in pool
		// Finds/counts entities in positive direction
		posDirIter := iterator.New(determTokensIter.Elems())
		posDirIter.SetPos(currDetermTokensPos)
		// Finds/counts entities in negative direction
		negDirIter := iterator.New(determTokensIter.Elems())
		negDirIter.SetPos(currDetermTokensPos)

		// [I, was, (with), Max, Payne, here] -> true, 2, Max Payne
		// [I, was, with, Max, Payne, (here)] -> false, 0, ""
		for posDirIter.Next() {
			// Should we include the entities here or substract it?
			entityDist++

			isEntity, entity := comp.TextWithEntities(posDirIter, entityTokensIter, comp.DirPos)
			if isEntity {
				appendMap(assocTokensAccum, determTokensIter, entityDist)
				// Skip about entity.
				posDirIter.Foward(len(entity) - 1) // Next increments
			}
		}

		// Reset distance
		entityDist = 0

		// [I, was, with, Max, Payne, (here)] -> true, 1, Max Payne
		// [I, was, (with), Max, Payne, here] -> false, 0, ""
		for negDirIter.Prev() {
			entityDist++

			isEntity, entity := comp.TextWithEntities(negDirIter, entityTokensIter, comp.DirNeg)
			if isEntity {
				appendMap(assocTokensAccum, determTokensIter, entityDist)
				negDirIter.Rewind(len(entity) - 1)
			}
		}
	}

	// Calculate the average distances
	for token, dist := range assocTokensAccum {
		assocTokens[token] = avg(dist)
	}
	return assocTokens, err
}

// Returns the average of a float slice
func avg(xs []float64) float64 {
	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}

// Helper to append float to a map
func appendMap(m map[string][]float64, k *iterator.Iterator[tokenize.Token], f float64) {
	text := k.CurrElem().Text
	m[text] = append(m[text], f)
}
