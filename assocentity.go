package assocentity

import (
	"context"
	"math"
	"strings"

	"github.com/ndabAP/assocentity/v12/internal/comp"
	"github.com/ndabAP/assocentity/v12/internal/iterator"
	"github.com/ndabAP/assocentity/v12/internal/pos"
	"github.com/ndabAP/assocentity/v12/tokenize"
)

// source wraps entities and texts
type source struct {
	Entities []string
	Texts    []string
}

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

	// Check if given PoS was found in text tokens
	if len(determTokens) == 0 {
		return dists, nil
	}

	// Create iterators

	determTokensIter := iterator.New(determTokens)

	// Search for entities in positive and negative direction
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
		for posDirIter.Next() {
			// [I, was, (with), Max, Payne, here] -> true, Max Payne
			// [I, was, with, Max, Payne, (here)] -> false, ""
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
			// [I, was, (with), Max, Payne, here] -> false, ""
			// [I, was, with, Max, Payne, (here)] -> true, Max Payne
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

// Aggregator aggregates tokens
type Aggregator func(tokenize.Token) tokenize.Token

// HumandReadableAggregator aggregates tokens through lower casing them and
// replacing them with their synonyms
var HumandReadableAggregator Aggregator = func(tok tokenize.Token) tokenize.Token {
	t := tokenize.Token{
		PoS:  tok.PoS,
		Text: strings.ToLower(tok.Text),
	}

	// This can increase the result data quality and could include more synonyms
	switch tok.Text {
	case "&":
		t.Text = "and"
	}

	return t
}

// Aggregate aggregates tokens with provided normalizer
func Aggregate(dists map[tokenize.Token][]float64, aggr Aggregator) {
	for tok, d := range dists {
		t := aggr(tok)

		// Check if text is the same as non-normalized
		if t == tok {
			continue
		}
		if _, ok := dists[t]; ok {
			dists[t] = append(dists[tok], d...)
		} else {
			dists[t] = d
		}

		delete(dists, tok)
	}
}

// Threshold exludes results that are below the given threshold. The threshold
// is described through the amount of distances per token relative to the total
// amount of tokens
func Threshold(dists map[tokenize.Token][]float64, threshold float64) {
	// Length of dists is amount of total tokens
	distsN := len(dists)
	for tok, d := range dists {
		dN := len(d)
		// Amount of distances per token relative to the amount of all tokens
		t := (float64(dN) / float64(distsN)) * 100
		if t < threshold {
			delete(dists, tok)
		}
	}
}
