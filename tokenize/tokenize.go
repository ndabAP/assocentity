package tokenize

import (
	"context"

	language "cloud.google.com/go/language/apiv1"
	option "google.golang.org/api/option"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

// Part of speech
const (
	ADJ   = 1 << iota // Adjective
	ADP               // Adposition
	ADV               // Adverb
	AFFIX             // Affix
	CONJ              // Conjunction
	DET               // Determiner
	NOUN              // Noun
	NUM               // Cardinal number
	PRON              // Pronoun
	PRT               // Particle or other function word
	PUNCT             // Punctuation
	UNKN              // Unknown
	VERB              // Verb (all tenses and modes)
	X                 // Other: foreign words, typos, abbreviations
	ANY   = ADJ | ADP | ADV | AFFIX | CONJ | DET | NOUN | NUM | PRON | PRT | PUNCT | UNKN | VERB | X
)

// Use map to be independent from library
var poS = map[languagepb.PartOfSpeech_Tag]int{
	languagepb.PartOfSpeech_ADJ:     ADJ,
	languagepb.PartOfSpeech_ADP:     ADP,
	languagepb.PartOfSpeech_ADV:     ADV,
	languagepb.PartOfSpeech_AFFIX:   AFFIX,
	languagepb.PartOfSpeech_CONJ:    CONJ,
	languagepb.PartOfSpeech_DET:     DET,
	languagepb.PartOfSpeech_NOUN:    NOUN,
	languagepb.PartOfSpeech_NUM:     NUM,
	languagepb.PartOfSpeech_PRON:    PRON,
	languagepb.PartOfSpeech_PRT:     PRT,
	languagepb.PartOfSpeech_PUNCT:   PUNCT,
	languagepb.PartOfSpeech_UNKNOWN: UNKN,
	languagepb.PartOfSpeech_VERB:    VERB,
	languagepb.PartOfSpeech_X:       X,
}

// Tokenizer tokenizes a text and entities
type Tokenizer interface {
	Tokenize(text string) ([]Token, error)
}

// Token represents a tokenized text unit
type Token struct {
	PoS   int    // Part of speech
	Token string // Text
}

var (
	client *language.Client
	err    error
	ctx    context.Context
)

// Lang defines the language used to examine the text. Both ISO and BCP-47 language codes are accepted
type Lang string

// AutoLang tries to automatically recognize the language
var AutoLang Lang = "auto"

// NLP tokenizes a text using NLP
type NLP struct {
	lang Lang
}

// NewNLP returns a new NLP instance
func NewNLP(credentialsFile string, lang Lang) (*NLP, error) {
	ctx = context.Background()

	// Create one client for all calls
	client, err = language.NewClient(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return &NLP{}, err
	}

	return &NLP{
		lang: lang,
	}, nil
}

// Tokenize tokenizes a text
func (nlp *NLP) Tokenize(text string) ([]Token, error) {
	var tokenized []Token
	tokenized, err = nlp.tokenize(text)
	if err != nil {
		return []Token{}, err
	}

	return tokenized, nil
}

// tokenize does the actual tokenization work
func (nlp *NLP) tokenize(text string) ([]Token, error) {
	resp, err := nlp.req(text)
	if err != nil {
		return nil, err
	}

	// Holds the tokenized text
	var tokenized []Token
	for _, t := range resp.GetTokens() {
		tokenized = append(tokenized, Token{
			PoS:   poS[t.PartOfSpeech.Tag],
			Token: t.GetText().GetContent(),
		})
	}

	return tokenized, nil
}

// req sends a request to the Google NLP server
func (nlp *NLP) req(text string) (*languagepb.AnnotateTextResponse, error) {
	doc := &languagepb.Document{
		Source: &languagepb.Document_Content{
			Content: text,
		},
		Type: languagepb.Document_PLAIN_TEXT,
	}

	if nlp.lang != "auto" {
		doc.Language = string(nlp.lang)
	}

	return client.AnnotateText(ctx, &languagepb.AnnotateTextRequest{
		Document: doc,
		Features: &languagepb.AnnotateTextRequest_Features{
			ExtractSyntax: true,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
}
