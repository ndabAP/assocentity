package tokenize

import (
	"context"

	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

// Tokenizer tokenizes a text
type Tokenizer interface {
	TokenizeText() ([]string, error)
	TokenizeEntities() ([][]string, error)
}

// NLP tokenizes a text using NLP
type NLP struct {
	text     string
	entities []string
	punct    bool // Punctation
}

// NewNLP returns a new NLP instance
func NewNLP(text string, entities []string, punct bool) NLP {
	return NLP{
		text:     text,
		entities: entities,
		punct:    punct,
	}
}

// Tokenize tokenizes a text
func (nlp *NLP) Tokenize() ([]string, error) {
	return tokenize(nlp.text, nlp.punct)
}

// TokenizedNested returns nested tokenized entities
func (nlp *NLP) TokenizedNested() ([][]string, error) {
	var tokenizedEntities [][]string
	for idx, entity := range nlp.entities {
		tokenizedEntity, err := tokenize(entity, nlp.punct)
		if err != nil {
			return nil, err
		}

		tokenizedEntities[idx] = tokenizedEntity
	}

	return tokenizedEntities, nil
}

func tokenize(text string, punct bool) ([]string, error) {
	ctx := context.Background()

	// Create a client
	client, err := language.NewClient(ctx)
	defer client.Close()
	if err != nil {
		return nil, err
	}

	resp, err := client.AnnotateText(ctx, &languagepb.AnnotateTextRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		Features: &languagepb.AnnotateTextRequest_Features{
			ExtractSyntax: true,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})

	if err != nil {
		return nil, err
	}

	// Holds the tokenized text
	var tokenizedText []string
	for _, v := range resp.GetTokens() {
		// Check for punctation
		if (v.PartOfSpeech.Tag != languagepb.PartOfSpeech_PUNCT) && !punct {
			tokenizedText = append(tokenizedText, v.GetText().GetContent())
		}
	}

	return tokenizedText, nil
}
