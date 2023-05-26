package nlp

import (
	"context"
	"errors"
	"time"

	language "cloud.google.com/go/language/apiv1"
	"github.com/googleapis/gax-go/v2/apierror"
	"github.com/ndabAP/assocentity/v14/tokenize"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/api/error_reason"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

var (
	ErrMaxRetries = errors.New("max retries reached")
)

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

// AutoLang tries to automatically recognize the language
var AutoLang string = "auto"

// NLPTokenizer tokenizes a text using Google NLP
type NLPTokenizer struct {
	credsFilename string
	lang          string
}

// NewNLPTokenizer returns a new NLP tokenizer instance. Note that NLPTokenizer
// has a built-in retrier
func NewNLPTokenizer(credentialsFilename string, lang string) tokenize.Tokenizer {
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
		if _, ok := poSMap[tok.PartOfSpeech.Tag]; !ok {
			return tokens, errors.New("can't find pos match")
		}

		tokens = append(tokens, tokenize.Token{
			PoS:  poSMap[tok.PartOfSpeech.Tag],
			Text: tok.GetText().GetContent(),
		})
	}
	return tokens, nil
}

// req sends a request to the Google server. It retries if the API rate limited
// is reached
func (nlp NLPTokenizer) req(ctx context.Context, text string) (*languagepb.AnnotateTextResponse, error) {
	client, err := language.NewClient(ctx, option.WithCredentialsFile(nlp.credsFilename))
	if err != nil {
		return &languagepb.AnnotateTextResponse{}, err
	}

	defer client.Close()

	doc := &languagepb.Document{
		Source: &languagepb.Document_Content{
			Content: text,
		},
		Type: languagepb.Document_PLAIN_TEXT,
	}
	// Set the desired language if not auto
	if nlp.lang != AutoLang {
		doc.Language = nlp.lang
	}

	// Google rate limit timeout
	const apiRateTimeout = 1.0 // In Minutes
	var (
		// Google errors
		apiErr                     *apierror.APIError
		errReasonRateLimitExceeded = error_reason.ErrorReason_RATE_LIMIT_EXCEEDED.String()

		delay     = apiRateTimeout
		delayMult = 1.05 // Delay multiplier
		retries   = 0
	)
	const (
		delayGrowth = 1.05 // Delay growth rate
		maxRetries  = 6
	)
	// Retry request up to maxRetries times if rate limit exceeded with an
	// growing delay
	for {
		if retries >= maxRetries {
			return &languagepb.AnnotateTextResponse{}, ErrMaxRetries
		}

		// Do the actual request
		res, err := client.AnnotateText(ctx, &languagepb.AnnotateTextRequest{
			Document: doc,
			Features: &languagepb.AnnotateTextRequest_Features{
				ExtractSyntax: true,
			},
			EncodingType: languagepb.EncodingType_UTF8,
		})
		// Check for rate limit exceeded error to retry
		if errors.As(err, &apiErr) {
			if apiErr.Reason() == errReasonRateLimitExceeded {
				time.Sleep(time.Minute * time.Duration(delay))

				// Retryer logic
				retries += 1
				delay *= delayMult
				delayMult *= delayGrowth

				continue
			}
		} else {
			return res, err
		}
	}
}
