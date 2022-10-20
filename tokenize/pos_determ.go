package tokenize

// PoSDetermer determinates if part of speech tags should be deleted
type PoSDetermer interface {
	DetermPoS(textTokens []Token, entitiesTokens [][]Token) ([]Token, error)
}
