package assocentity

import (
	"strings"

	"github.com/ndabAP/assocentity/v13/tokenize"
)

// Normalizer normalizes tokens like lower casing them
type Normalizer func(tokenize.Token) tokenize.Token

// HumandReadableNormalizer normalizes tokens through lower casing them and
// replacing them with their synonyms
var HumandReadableNormalizer Normalizer = func(tok tokenize.Token) tokenize.Token {
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

// Normalize normalizes tokens with provided normalizer
func Normalize(dists map[tokenize.Token][]float64, norm Normalizer) {
	for tok, d := range dists {
		t := norm(tok)

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

// Threshold excludes results that are below the given threshold. The threshold
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

var EnglishDictonary = map[tokenize.PoS]func(tokenize.Token) bool{}

func Dictonary(dists map[tokenize.Token][]float64, dict map[tokenize.PoS]func(tokenize.Token) bool) {
	for tok := range dists {
		f, ok := dict[tok.PoS]
		if !ok {
			continue
		}

		if !f(tok) {
			delete(dists, tok)
		}
	}
}

func LoadDict() error {
	return nil
}
