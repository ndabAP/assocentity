package tokenize

// PoSDetermer determinates which part of speech tags should be kept. It
// receives the tokenized text and tokenized entities and returns the tokenized
// text while only the desired part of speeches are kept. Entities must be
// always kept
type PoSDetermer interface {
	DetermPoS(textTokens []Token, entitiesTokens [][]Token) []Token
}
