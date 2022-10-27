package assocentity

import (
	"context"

	"github.com/ndabAP/assocentity/v8/internal/comp"
	"github.com/ndabAP/assocentity/v8/internal/iterator"
	"github.com/ndabAP/assocentity/v8/tokenize"
)

// Algorithm:
//
// 	1. Tokenize text: [Without, Mona, 's, help]
//	2. Tokenize entities: [[Max], [Max, Payne]]
//	3. Remove part of speech:
//	   [Vlad, was, right, ., There, are, no, choices] becomes without Verbs
//	   [Vlad, right, ., There,  no, choices]
// 	4. Iterate through part of speech removed text tokens
//	4. a) If the current text token and further equals an entity, skip text
//		  about entity positions
//	   b) Iterate in positive direction through part of speech removed text
//		  tokens
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
		// If the current token is an entity, we skip about the entity
		currDetermTokensPos := determTokensIter.CurrPos()
		isEntity, entity := comp.TextWithEntity(determTokensIter, entityTokensIter, comp.DirPos)
		if isEntity {
			determTokensIter.SetPos(currDetermTokensPos + len(entity))
		}

		// Now we can collect the actual distances

		// Distance
		var entityDist float64

		// Finds/counts entities in positive direction
		posDirIter := determTokensIter
		posDirIter.Reset().SetPos(currDetermTokensPos)
		// Finds/counts entities in negative direction
		negDirIter := determTokensIter
		negDirIter.Reset().SetPos(currDetermTokensPos)

		// [I, was, (with), Max, Payne, here] -> true, 2, Max Payne
		// [I, was, with, Max, Payne, (here)] -> false, 0, ""
		for posDirIter.Next() {
			currPosDirPos := negDirIter.CurrPos()
			isEntity, entity := comp.TextWithEntity(posDirIter, entityTokensIter, comp.DirPos)
			if isEntity {
				entityDist++

				appendMap(assocEntitiesVals, posDirIter.CurrElem().Text, entityDist)
				// Skip about entity
				posDirIter.SetPos(currPosDirPos + len(entity))
			}

		}

		// Reset distance
		entityDist = 0

		// [I, was, with, Max, Payne, (here)] -> true, 1, Max Payne
		// [I, was, (with), Max, Payne, here] -> false, 0, ""
		for negDirIter.Prev() {
			currNegDirPos := negDirIter.CurrPos()
			isEntity, entity := comp.TextWithEntity(posDirIter, entityTokensIter, comp.DirNeg)
			if isEntity {
				entityDist++

				appendMap(assocEntitiesVals, negDirIter.CurrElem().Text, entityDist)
				posDirIter.SetPos(currNegDirPos - len(entity))
			}

		}
	}

	// Calculate the average distances
	for token, dist := range assocEntitiesVals {
		assocEntities[token] = avg(dist)
	}
	return assocEntities, err
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
func appendMap(m map[string][]float64, k string, f float64) {
	m[k] = append(m[k], f)
}
