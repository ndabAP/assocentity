package tokenize

import (
	"context"
)

// Part of speech
type PoS int

const (
	ANY = ADJ | ADP | ADV | AFFIX | CONJ | DET | NOUN | NUM | PRON | PRT | PUNCT | UNKN | VERB | X

	UNKN PoS = 1 << iota // Unknown
	X                    // Other: foreign words, typos, abbreviations

	ADJ   // Adjective
	ADP   // Adposition
	ADV   // Adverb
	AFFIX // Affix
	CONJ  // Conjunction
	DET   // Determiner
	NOUN  // Noun
	NUM   // Cardinal number
	PRON  // Pronoun
	PRT   // Particle or other function word
	PUNCT // Punctuation
	VERB  // Verb (all tenses and modes)
)

// Tokenizer tokenizes a text and entities
type Tokenizer interface {
	Tokenize(ctx context.Context, text string) ([]Token, error)
}

// Token represents a tokenized text unit
type Token struct {
	PoS  PoS    // Part of speech
	Text string // Text
}
