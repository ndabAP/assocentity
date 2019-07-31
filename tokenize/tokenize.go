package tokenize

import (
	"context"

	language "cloud.google.com/go/language/apiv1"
	option "google.golang.org/api/option"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

// Tokenizer tokenizes a text
type Tokenizer interface {
	TokenizeText() ([]string, error)
	TokenizeEntities() ([][]string, error)
}

var client *language.Client
var err error
var ctx context.Context

// NLP tokenizes a text using NLP
type NLP struct {
	text     string
	entities []string
	punct    bool // Punctation
}

// NewNLP returns a new NLP instance
func NewNLP(credentialsFile, text string, entities []string, punct bool) (*NLP, error) {
	ctx = context.Background()

	// Create one client for all calls
	client, err = language.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return &NLP{}, err
	}

	return &NLP{
		text:     text,
		entities: entities,
		punct:    punct,
	}, nil
}

// Tokenize tokenizes a text
func (nlp *NLP) Tokenize() ([]string, error) {
	return nlp.tokenize(nlp.text, nlp.punct)
}

// TokenizedNested returns nested tokenized entities
func (nlp *NLP) TokenizedNested() ([][]string, error) {
	var tokenizedEntities [][]string
	for idx, entity := range nlp.entities {
		tokenizedEntity, err := nlp.tokenize(entity, nlp.punct)
		if err != nil {
			return nil, err
		}

		tokenizedEntities[idx] = tokenizedEntity
	}

	return tokenizedEntities, nil
}

// tokenize does the actual tokenization work
func (nlp *NLP) tokenize(text string, punct bool) ([]string, error) {
	resp, err := client.AnnotateText(ctx, &languagepb.AnnotateTextRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: nlp.text,
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
		if nlp.punct {
			tokenizedText = append(tokenizedText, v.GetText().GetContent())
		} else {
			if v.PartOfSpeech.Tag != languagepb.PartOfSpeech_PUNCT {
				tokenizedText = append(tokenizedText, v.GetText().GetContent())
			}
		}
	}

	return tokenizedText, nil
}
