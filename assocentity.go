package assocentity

import (
	"context"
	"math"

	"github.com/ndabAP/assocentity/v14/internal/comp"
	"github.com/ndabAP/assocentity/v14/internal/iterator"
	"github.com/ndabAP/assocentity/v14/internal/pos"
	"github.com/ndabAP/assocentity/v14/tokenize"
)

// source wraps entities and texts
type source struct {
	Entities []string
	Texts    []string
}

// NewSource returns a new source consisting of entities and texts
func NewSource(entities, texts []string) source {
	return source{
		Entities: entities,
		Texts:    texts,
	}
}

// Distances returns the distances from entities to a list of texts
func Distances(
	ctx context.Context,
	tokenizer tokenize.Tokenizer,
	poS tokenize.PoS,
	source source,
) (map[tokenize.Token][]float64, error) {
	var (
		dists = make(map[tokenize.Token][]float64)
		err   error
	)
	for _, text := range source.Texts {
		d, err := distances(ctx, tokenizer, poS, text, source.Entities)
		if err != nil {
			return dists, err
		}

		for tok, dist := range d {
			dists[tok] = append(dists[tok], dist...)
		}
	}

	return dists, err
}

// distances returns the distances to entities for one text
func distances(
	ctx context.Context,
	tokenizer tokenize.Tokenizer,
	poS tokenize.PoS,
	text string,
	entities []string,
) (map[tokenize.Token][]float64, error) {
	var (
		dists = make(map[tokenize.Token][]float64)
		err   error
	)

	// Tokenize text
	textTokens, err := tokenizer.Tokenize(ctx, text)
	if err != nil {
		return dists, err
	}

	// Tokenize entities
	var entityTokens [][]tokenize.Token
	for _, entity := range entities {
		tokens, err := tokenizer.Tokenize(ctx, entity)
		if err != nil {
			return dists, err
		}
		entityTokens = append(entityTokens, tokens)
	}

	// Determinate part of speech
	posDetermer := pos.NewPoSDetermer(poS)
	determTokens := posDetermer.DetermPoS(textTokens, entityTokens)

	// Check if any given PoS was found in text tokens
	if len(determTokens) == 0 {
		return dists, nil
	}

	// Create iterators

	determTokensIter := iterator.New(determTokens)

	// Use iterators to search for entities in positive and negative direction
	posDirIter := iterator.New(determTokens)
	negDirIter := iterator.New(determTokens)

	entityTokensIter := iterator.New(entityTokens)

	// Iterate through part of speech determinated text tokens
	for determTokensIter.Next() {
		// If the current text token is an entity, we skip about the entity
		currDetermTokensPos := determTokensIter.CurrPos()
		isEntity, entity := comp.TextWithEntities(
			determTokensIter,
			entityTokensIter,
			comp.DirPos,
		)
		if isEntity {
			determTokensIter.Forward(len(entity) - 1)
			continue
		}

		// Now we can collect the actual distances

		// Finds/counts entities in positive direction
		posDirIter.SetPos(currDetermTokensPos)
		for posDirIter.Next() {
			isEntity, entity := comp.TextWithEntities(
				posDirIter,
				entityTokensIter,
				comp.DirPos,
			)
			if isEntity {
				appendDist(dists, determTokensIter, posDirIter)
				// Skip about entity
				posDirIter.Forward(len(entity) - 1) // Next increments
			}
		}

		// Finds/counts entities in negative direction
		negDirIter.SetPos(currDetermTokensPos)
		for negDirIter.Prev() {
			isEntity, entity := comp.TextWithEntities(
				negDirIter,
				entityTokensIter,
				comp.DirNeg,
			)
			if isEntity {
				appendDist(dists, determTokensIter, negDirIter)
				negDirIter.Rewind(len(entity) - 1)
			}
		}
	}

	return dists, err
}

// Helper to append a float64 to a map of tokens and distances
func appendDist(
	m map[tokenize.Token][]float64,
	k *iterator.Iterator[tokenize.Token],
	v *iterator.Iterator[tokenize.Token],
) {
	token := k.CurrElem()
	dist := math.Abs(float64(v.CurrPos() - k.CurrPos()))
	m[token] = append(m[token], dist)
}

// Mean returns the mean of the provided distances
func Mean(dists map[tokenize.Token][]float64) map[tokenize.Token]float64 {
	mean := make(map[tokenize.Token]float64)
	for token, d := range dists {
		mean[token] = meanFloat64(d)
	}
	return mean
}

// Returns the mean of a 64-bit float slice
func meanFloat64(xs []float64) float64 {
	// Prevent /0
	if len(xs) == 0 {
		return 0
	}

	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}
