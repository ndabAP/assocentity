package assocentity_test

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	"github.com/ndabAP/assocentity/v9"
	"github.com/ndabAP/assocentity/v9/nlp"
	"github.com/ndabAP/assocentity/v9/tokenize"
)

func TestDoSimple1(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	credentialsFile := os.Getenv("GOOGLE_NLP_SERVICE_ACCOUNT_FILE_LOCATION")
	nlpTokenizer := nlp.NewNLPTokenizer(credentialsFile, nlp.AutoLang)

	text := "Relax, Max. You're a nice guy."
	entities := []string{"Max", "Max Payne"}

	posDeterm := nlp.NewNLPPoSDetermer(tokenize.ANY)

	got, err := assocentity.Do(context.Background(), nlpTokenizer, posDeterm, text, entities)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]float64{
		"Relax": 2,
		",":     1,
		".":     4,
		"You":   2,
		"'re":   3,
		"a":     4,
		"nice":  5,
		"guy":   6,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Do() = %v, want %v", got, want)
	}
}

func TestDoSimple2(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	credentialsFile := os.Getenv("GOOGLE_NLP_SERVICE_ACCOUNT_FILE_LOCATION")
	nlpTokenizer := nlp.NewNLPTokenizer(credentialsFile, nlp.AutoLang)

	text := "Punchinello wanted Payne? He'd see the pain."
	entities := []string{"Punchinello", "Payne"}

	dps := nlp.NewNLPPoSDetermer(tokenize.ANY)

	got, err := assocentity.Do(context.Background(), nlpTokenizer, dps, text, entities)
	if err != nil {
		log.Fatal(err)
	}

	want := map[string]float64{
		"wanted": 1,
		"?":      2,
		"He":     3,
		"'d":     4,
		"see":    5,
		"the":    6,
		"pain":   7,
		".":      8,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Assoc() = %v, want %v", got, want)
	}
}

func TestDoComplex(t *testing.T) {
	text := "ee ee aa bb cc dd. b ff, gg, hh, bb, bb. ii!"
	entities := []string{"bb", "b", "ee ee"}

	posDeterm := nlp.NewNLPPoSDetermer(tokenize.ANY)

	var tTokenizer testTokenizer
	got, err := assocentity.Do(context.Background(), tTokenizer, posDeterm, text, entities)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]float64{
		"aa": 6.6,
		",":  6.3,
		"!":  10.8,
		"cc": 5.8,
		"dd": 5.6,
		".":  7.1,
		"ff": 5.4,
		"gg": 5.8,
		"hh": 6.2,
		"ii": 9.8,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Do() = %v, want %v", got, want)
	}
}

// Mock tokenization
type testTokenizer int

// Hack to simulate different tokenization response steps
var tokCall int

// Mock date: 10-30-2022
func (tt testTokenizer) Tokenize(ctx context.Context, text string) ([]tokenize.Token, error) {
	tokCall++

	switch tokCall {
	case 1:
		return []tokenize.Token{
			{PoS: tokenize.NOUN, Text: "ee"},
			{PoS: tokenize.NOUN, Text: "ee"},
			{PoS: tokenize.NOUN, Text: "aa"},
			{PoS: tokenize.NOUN, Text: "bb"},
			{PoS: tokenize.NOUN, Text: "cc"},
			{PoS: tokenize.NOUN, Text: "dd"},
			{PoS: tokenize.PUNCT, Text: "."},
			{PoS: tokenize.NOUN, Text: "b"},
			{PoS: tokenize.NOUN, Text: "ff"},
			{PoS: tokenize.PUNCT, Text: ","},
			{PoS: tokenize.X, Text: "gg"},
			{PoS: tokenize.PUNCT, Text: ","},
			{PoS: tokenize.NOUN, Text: "hh"},
			{PoS: tokenize.PUNCT, Text: ","},
			{PoS: tokenize.NOUN, Text: "bb"},
			{PoS: tokenize.PUNCT, Text: ","},
			{PoS: tokenize.NOUN, Text: "bb"},
			{PoS: tokenize.PUNCT, Text: "."},
			{PoS: tokenize.NOUN, Text: "ii"},
			{PoS: tokenize.PUNCT, Text: "!"},
		}, nil
	case 2:
		return []tokenize.Token{
			{PoS: tokenize.NOUN, Text: "ee"},
			{PoS: tokenize.NOUN, Text: "ee"},
		}, nil
	case 3:
		return []tokenize.Token{
			{PoS: tokenize.NOUN, Text: "bb"},
		}, nil
	case 4:
		return []tokenize.Token{
			{PoS: tokenize.NOUN, Text: "b"},
		}, nil
	case 5:
		return []tokenize.Token{
			{PoS: tokenize.NOUN, Text: "bb"},
		}, nil
	// 6
	default:
		return []tokenize.Token{
			{PoS: tokenize.NOUN, Text: "bb"},
		}, nil
	}
}
