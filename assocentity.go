package assocentity

import (
	"context"
	"math"

	"github.com/ndabAP/assocentity/v9/internal/comp"
	"github.com/ndabAP/assocentity/v9/internal/iterator"
	"github.com/ndabAP/assocentity/v9/tokenize"
)

// Do returns the average distance from entities to a text consisting of token
func Do(ctx context.Context, tokenizer tokenize.Tokenizer, psd tokenize.PoSDetermer, text string, entities []string) (map[string]float64, error) {
	var (
		assocEntities      = make(map[string]float64)
		assocEntitiesAccum = make(map[string][]float64)

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
		// If the current text token is an entity, we skip about the entity
		currDetermTokensPos := determTokensIter.CurrPos()
		isEntity, entity := comp.TextWithEntities(determTokensIter, entityTokensIter, comp.DirPos)
		if isEntity {
			determTokensIter.Forward(len(entity) - 1)
			continue
		}

		// Now we can collect the actual distances

		// Finds/counts entities in positive direction
		posDirIter := iterator.New(determTokensIter.Elems())
		posDirIter.SetPos(currDetermTokensPos)
		// Finds/counts entities in negative direction
		negDirIter := iterator.New(determTokensIter.Elems())
		negDirIter.SetPos(currDetermTokensPos)

		// [I, was, (with), Max, Payne, here] -> true, Max Payne
		// [I, was, with, Max, Payne, (here)] -> false, ""
		for posDirIter.Next() {
			isEntity, entity := comp.TextWithEntities(posDirIter, entityTokensIter, comp.DirPos)
			if isEntity {
				appendTokenDist(assocEntitiesAccum, determTokensIter, posDirIter)
				// Skip about entity.
				posDirIter.Forward(len(entity) - 1) // Next increments
			}
		}

		// [I, was, with, Max, Payne, (here)] -> true, Max Payne
		// [I, was, (with), Max, Payne, here] -> false, ""
		for negDirIter.Prev() {
			isEntity, entity := comp.TextWithEntities(negDirIter, entityTokensIter, comp.DirNeg)
			if isEntity {
				appendTokenDist(assocEntitiesAccum, determTokensIter, negDirIter)
				negDirIter.Rewind(len(entity) - 1)
			}
		}
	}

	// Calculate the average distances
	for token, dist := range assocEntitiesAccum {
		assocEntities[token] = avgFloat(dist)
	}
	return assocEntities, err
}

// Helper to append float to a map
func appendTokenDist(m map[string][]float64, k *iterator.Iterator[tokenize.Token], v *iterator.Iterator[tokenize.Token]) {
	token := k.CurrElem().Text
	dist := math.Abs(float64(v.CurrPos() - k.CurrPos()))
	m[token] = append(m[token], dist)
}

// Returns the average of a float slice
func avgFloat(xs []float64) float64 {
	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}
