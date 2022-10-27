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

	// Iterate through part of speech determinated text tokens to find the
	// entity
	for determTokensIter.Next() {
		// Skip about entity positions, if entity
		currDetermTokensPos := determTokensIter.CurrPos()
		isEntity, entity := comp.TextWithEntity(determTokensIter, entityTokensIter, comp.PosDir)
		if isEntity {
			determTokensIter.SetPos(currDetermTokensPos + len(entity))
		}

		// Distance
		var entityDist float64

		// Iterate in positive and negative direction to find entity (distances)
		// Finds/counts entities in positive direction
		posDirIter := determTokensIter
		// Finds/counts entities in negative direction
		negDirIter := determTokensIter

		// [I, was, (with), Max, Payne, here] -> true, 2, Max Payne
		// [I, was, with, Max, Payne, (here)]
		for posDirIter.Next() {
			currPosDirPos := negDirIter.CurrPos()
			isEntity, entity := comp.TextWithEntity(posDirIter, entityTokensIter, comp.PosDir)
			if isEntity {
				appendMap(assocEntitiesVals, posDirIter.CurrElem().Text, entityDist)
				// Skip about entity
				posDirIter.SetPos(currPosDirPos + len(entity))
			}

			entityDist++
		}

		// Reset distance
		entityDist = 0

		// [I, was, with, Max, Payne, (here)] -> true, 1, Max Payne
		// [I, was, (with), Max, Payne, here]
		for negDirIter.Prev() {
			currNegDirPos := negDirIter.CurrPos()
			isEntity, entity := comp.TextWithEntity(posDirIter, entityTokensIter, comp.NegDir)
			if isEntity {
				appendMap(assocEntitiesVals, negDirIter.CurrElem().Text, entityDist)
				posDirIter.SetPos(currNegDirPos - len(entity))
			}

			entityDist++
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
