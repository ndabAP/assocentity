package tokenize

// PoSDetermer determinates if part of speech tags should be deleted
type PoSDetermer interface {
	Determ(tokenizedText []Token, tokenizedEntities [][]Token) ([]Token, error)
}
