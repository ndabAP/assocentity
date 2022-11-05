package assocentity_test

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	"github.com/ndabAP/assocentity/v10"
	"github.com/ndabAP/assocentity/v10/nlp"
	"github.com/ndabAP/assocentity/v10/tokenize"
)

func TestDoWired(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	credentialsFile := os.Getenv("GOOGLE_NLP_SERVICE_ACCOUNT_FILE_LOCATION")
	nlpTokenizer := nlp.NewNLPTokenizer(credentialsFile, nlp.AutoLang)

	t.Run("rand1", func(t *testing.T) {
		text := "Relax, Max. You're a nice guy."
		entities := []string{"Max", "Max Payne"}

		posDeterm := nlp.NewNLPPoSDetermer(tokenize.ANY)

		got, err := assocentity.Do(context.Background(), nlpTokenizer, posDeterm, text, entities)
		if err != nil {
			t.Fatal(err)
		}
		want := map[tokenize.Token]float64{
			{
				PoS:  tokenize.VERB,
				Text: "Relax",
			}: 2,
			{
				PoS:  tokenize.PUNCT,
				Text: ",",
			}: 1,
			{
				PoS:  tokenize.PUNCT,
				Text: ".",
			}: 4,
			{
				PoS:  tokenize.PRON,
				Text: "You",
			}: 2,
			{
				PoS:  tokenize.VERB,
				Text: "'re",
			}: 3,
			{
				PoS:  tokenize.DET,
				Text: "a",
			}: 4,
			{
				PoS:  tokenize.ADJ,
				Text: "nice",
			}: 5,
			{
				PoS:  tokenize.NOUN,
				Text: "guy",
			}: 6,
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Do() = %v, want %v", got, want)
		}
	})

	t.Run("rand2", func(t *testing.T) {
		text := "Punchinello wanted Payne? He'd see the pain."
		entities := []string{"Punchinello", "Payne"}

		dps := nlp.NewNLPPoSDetermer(tokenize.ANY)

		got, err := assocentity.Do(context.Background(), nlpTokenizer, dps, text, entities)
		if err != nil {
			log.Fatal(err)
		}

		want := map[tokenize.Token]float64{
			{
				PoS:  tokenize.VERB,
				Text: "wanted",
			}: 1,
			{
				PoS:  tokenize.PUNCT,
				Text: "?",
			}: 2,
			{
				PoS:  tokenize.PRON,
				Text: "He",
			}: 3,
			{
				PoS:  tokenize.VERB,
				Text: "'d",
			}: 4,
			{
				PoS:  tokenize.VERB,
				Text: "see",
			}: 5,
			{
				PoS:  tokenize.DET,
				Text: "the",
			}: 6,
			{
				PoS:  tokenize.NOUN,
				Text: "pain",
			}: 7,
			{
				PoS:  tokenize.PUNCT,
				Text: ".",
			}: 8,
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Assoc() = %v, want %v", got, want)
		}
	})
}

func TestDoWireless(t *testing.T) {
	text := "ee ee aa bb cc dd. b ff, gg, hh, bb, bb. ii!"
	entities := []string{"bb", "b", "ee ee"}

	posDeterm := nlp.NewNLPPoSDetermer(tokenize.ANY)

	var tTokenizer testTokenizer
	got, err := assocentity.Do(context.Background(), tTokenizer, posDeterm, text, entities)
	if err != nil {
		t.Fatal(err)
	}
	want := map[tokenize.Token]float64{
		{
			PoS:  tokenize.NOUN,
			Text: "aa",
		}: 6.6,
		{
			PoS:  tokenize.PUNCT,
			Text: ",",
		}: 6.3,
		{
			PoS:  tokenize.PUNCT,
			Text: "!",
		}: 10.8,
		{
			PoS:  tokenize.NOUN,
			Text: "cc",
		}: 5.8,
		{
			PoS:  tokenize.NOUN,
			Text: "dd",
		}: 5.6,
		{
			PoS:  tokenize.PUNCT,
			Text: ".",
		}: 7.1,
		{
			PoS:  tokenize.NOUN,
			Text: "ff",
		}: 5.4,
		{
			PoS:  tokenize.X,
			Text: "gg",
		}: 5.8,
		{
			PoS:  tokenize.NOUN,
			Text: "hh",
		}: 6.2,
		{
			PoS:  tokenize.NOUN,
			Text: "ii",
		}: 9.8,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Do() = %v, want %v", got, want)
	}
}

// Mock tokenization
type testTokenizer int

// Hack to simulate different tokenization response steps
var calls int

// Mock date: 10-30-2022
func (tt testTokenizer) Tokenize(ctx context.Context, text string) ([]tokenize.Token, error) {
	calls++

	switch calls {
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
