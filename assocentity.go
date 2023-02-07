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

// Source wraps entities and texts
type Source struct {
	Entities []string
	Texts    []string
}

// Dists returns the distances from entities to a slice of texts. It ignores
// empty texts and not found pos
func Dists(
	ctx context.Context,
	tokenizer tokenize.Tokenizer,
	poS tokenize.PoS,
	source Source,
) (map[[2]string][]float64, error) {
	var (
		dists = make(map[[2]string][]float64)
		err   error
	)

	for _, text := range source.Texts {
		d, err := distances(ctx, tokenizer, poS, text, source.Entities)
		if err != nil {
			return dists, err
		}

		for token, dist := range d {
			t := [2]string{token.Text, tokenize.PoSMapStr[token.PoS]}
			dists[t] = append(dists[t], dist...)
		}
	}

	return dists, err
}

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
		// [I, was, (with), Max, Payne, here] -> true, Max Payne
		// [I, was, with, Max, Payne, (here)] -> false, ""
		for posDirIter.Next() {
			isEntity, entity := comp.TextWithEntities(posDirIter, entityTokensIter, comp.DirPos)
			if isEntity {
				appendTokenDist(dists, determTokensIter, posDirIter)
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
				appendTokenDist(dists, determTokensIter, negDirIter)
				negDirIter.Rewind(len(entity) - 1)
			}
		}
	}

	return dists, err
}

// Helper to append float to a map
func appendTokenDist(m map[tokenize.Token][]float64, k *iterator.Iterator[tokenize.Token], v *iterator.Iterator[tokenize.Token]) {
	token := k.CurrElem()
	dist := math.Abs(float64(v.CurrPos() - k.CurrPos()))
	m[token] = append(m[token], dist)
}

func Count(dists map[[2]string][]float64) [][2]string {
	c := make([][2]string, len(dists))
	return c
}

func Mean(dists map[[2]string][]float64) map[[2]string]float64 {
	mean := make(map[[2]string]float64)
	for token, d := range dists {
		mean[token] = meanFloat64(d)
	}
	return mean
}

// Returns the mean of a 64-bit float slice
func meanFloat64(xs []float64) float64 {
	sum := 0.0
	for _, x := range xs {
		sum += x
	}
	return sum / float64(len(xs))
}

func Normalize(dists map[[2]string][]float64) {
	for tok, d := range dists {
		t := strings.ToLower(tok[1])
		dists[[2]string{tok[0], t}] = d
	}
}

func Threshold(dists map[[2]string][]float64, threshold int) {
	distsN := len(dists)
	for tok, tokDist := range dists {
		tokDistN := len(tokDist)
		if tokDistN/distsN < threshold {
			delete(dists, tok)
		}
	}
}
