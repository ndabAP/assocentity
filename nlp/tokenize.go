package nlp

import (
	"context"

	language "cloud.google.com/go/language/apiv1"
	"github.com/ndabAP/assocentity/v9/tokenize"
	"google.golang.org/api/option"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

// Use map to be independent from library
var poSMap = map[languagepb.PartOfSpeech_Tag]tokenize.PoS{
	languagepb.PartOfSpeech_ADJ:     tokenize.ADJ,
	languagepb.PartOfSpeech_ADP:     tokenize.ADP,
	languagepb.PartOfSpeech_ADV:     tokenize.ADV,
	languagepb.PartOfSpeech_AFFIX:   tokenize.AFFIX,
	languagepb.PartOfSpeech_CONJ:    tokenize.CONJ,
	languagepb.PartOfSpeech_DET:     tokenize.DET,
	languagepb.PartOfSpeech_NOUN:    tokenize.NOUN,
	languagepb.PartOfSpeech_NUM:     tokenize.NUM,
	languagepb.PartOfSpeech_PRON:    tokenize.PRON,
	languagepb.PartOfSpeech_PRT:     tokenize.PRT,
	languagepb.PartOfSpeech_PUNCT:   tokenize.PUNCT,
	languagepb.PartOfSpeech_UNKNOWN: tokenize.UNKN,
	languagepb.PartOfSpeech_VERB:    tokenize.VERB,
	languagepb.PartOfSpeech_X:       tokenize.X,
}

// Lang defines the language used to examine the text. Both ISO and BCP-47
// language codes are accepted
type Lang string

// AutoLang tries to automatically recognize the language
var AutoLang Lang = "auto"

// NLPTokenizer tokenizes a text using Google NLP
type NLPTokenizer struct {
	credsFilename string
	lang          Lang
}

// NewNLPTokenizer returns a new NLP tokenizer instance
func NewNLPTokenizer(credentialsFilename string, lang Lang) tokenize.Tokenizer {
	return NLPTokenizer{
		credsFilename: credentialsFilename,
		lang:          lang,
	}
}

// Tokenize tokenizes a text
func (nlp NLPTokenizer) Tokenize(ctx context.Context, text string) ([]tokenize.Token, error) {
	res, err := nlp.req(ctx, text)
	if err != nil {
		return []tokenize.Token{}, err
	}

	tokens := make([]tokenize.Token, 0)
	for _, tok := range res.GetTokens() {
		tokens = append(tokens, tokenize.Token{
			PoS:  poSMap[tok.PartOfSpeech.Tag],
			Text: tok.GetText().GetContent(),
		})
	}
	return tokens, nil
}

// req sends a request to the Google server
func (nlp NLPTokenizer) req(ctx context.Context, text string) (*languagepb.AnnotateTextResponse, error) {
	client, err := language.NewClient(ctx, option.WithCredentialsFile(nlp.credsFilename))
	if err != nil {
		return &languagepb.AnnotateTextResponse{}, err
	}

	doc := &languagepb.Document{
		Source: &languagepb.Document_Content{
			Content: text,
		},
		Type: languagepb.Document_PLAIN_TEXT,
	}
	// Set desired language if not auto
	if nlp.lang != AutoLang {
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
