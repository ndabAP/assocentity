package assocentity

import (
	"context"
	"math"

	"github.com/ndabAP/assocentity/v11/internal/comp"
	"github.com/ndabAP/assocentity/v11/internal/iterator"
	"github.com/ndabAP/assocentity/v11/internal/pos"
	"github.com/ndabAP/assocentity/v11/tokenize"
)

func MeanN(
	ctx context.Context,
	tokenizer tokenize.Tokenizer,
	poS tokenize.PoS,
	texts []string,
	entities []string,
) (map[tokenize.Token]float64, error) {
	mean := make(map[tokenize.Token]float64)

	means := make(map[tokenize.Token][]float64)
	for _, text := range texts {
		dists, err := dist(ctx, tokenizer, poS, text, entities)
		if err != nil {
			return mean, err
		}

		for tok, dist := range dists {
			means[tok] = append(means[tok], dist...)
		}
	}

	// Calculate the average distances
	for token, dist := range means {
		mean[token] = meanFloat64(dist)
	}
	return mean, nil
}

// Mean returns the average distance from entities to a text consisting of token
func Mean(
	ctx context.Context,
	tokenizer tokenize.Tokenizer,
	poS tokenize.PoS,
	text string,
	entities []string,
) (map[tokenize.Token]float64, error) {
	mean := make(map[tokenize.Token]float64)

	dists, err := dist(ctx, tokenizer, poS, text, entities)
	if err != nil {
		return mean, err
	}
	// Calculate the average distances
	for token, dist := range dists {
		mean[token] = meanFloat64(dist)
	}
	return mean, err
}

// func CountA amount of hits

// func TopN top n closed tokens

// func Closest top 1 most closest frequent

func dist(
	ctx context.Context,
	tokenizer tokenize.Tokenizer,
	poS tokenize.PoS,
	text string,
	entities []string,
) (map[tokenize.Token][]float64, error) {
	var (
		dist = make(map[tokenize.Token][]float64)
		err  error
	)

	// Tokenize text
	textTokens, err := tokenizer.Tokenize(ctx, text)
	if err != nil {
		return dist, err
	}

	// Tokenize entities
	var entityTokens [][]tokenize.Token
	for _, entity := range entities {
		tokens, err := tokenizer.Tokenize(ctx, entity)
		if err != nil {
			return dist, err
		}
		entityTokens = append(entityTokens, tokens)
	}

	// Determinate part of speech
	posDetermer := pos.NewPoSDetermer(poS)
	determTokens := posDetermer.DetermPoS(textTokens, entityTokens)

	// Creates iterators

	determTokensIter := iterator.New(determTokens)

	// Iterators to search for entities in positive and negative direction
	posDirIter := iterator.New(determTokens)
	negDirIter := iterator.New(determTokens)

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
		posDirIter.SetPos(currDetermTokensPos)
		// [I, was, (with), Max, Payne, here] -> true, Max Payne
		// [I, was, with, Max, Payne, (here)] -> false, ""
		for posDirIter.Next() {
			isEntity, entity := comp.TextWithEntities(posDirIter, entityTokensIter, comp.DirPos)
			if isEntity {
				appendTokenDist(dist, determTokensIter, posDirIter)
				// Skip about entity
				posDirIter.Forward(len(entity) - 1) // Next increments
			}
		}

		// Finds/counts entities in negative direction
		negDirIter.SetPos(currDetermTokensPos)
		// [I, was, with, Max, Payne, (here)] -> true, Max Payne
		// [I, was, (with), Max, Payne, here] -> false, ""
		for negDirIter.Prev() {
			isEntity, entity := comp.TextWithEntities(negDirIter, entityTokensIter, comp.DirNeg)
			if isEntity {
				appendTokenDist(dist, determTokensIter, negDirIter)
				negDirIter.Rewind(len(entity) - 1)
			}
		}
	}

	return dist, err
}

// Helper to append float to a map
func appendTokenDist(m map[tokenize.Token][]float64, k *iterator.Iterator[tokenize.Token], v *iterator.Iterator[tokenize.Token]) {
	token := k.CurrElem()
	dist := math.Abs(float64(v.CurrPos() - k.CurrPos()))
	m[token] = append(m[token], dist)
}

// Returns the average of a float slice
func meanFloat64(xs []float64) float64 {
	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}
