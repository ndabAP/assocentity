package assocentity_test

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	"github.com/ndabAP/assocentity/v8"
	"github.com/ndabAP/assocentity/v8/nlp"
	"github.com/ndabAP/assocentity/v8/tokenize"
)

var nlpTokenizer nlp.NLPTokenizer

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	credentialsFile := os.Getenv("GOOGLE_NLP_SERVICE_ACCOUNT_FILE_LOCATION")
	nlpTokenizer = nlp.NewNLPTokenizer(credentialsFile, nlp.AutoLang)
}

func TestAssocIntegrationSingleWordEntities(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

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
		".":     3.5,
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

// func TestAssocIntegrationSingleWordEntitiesEnglishLanguage(t *testing.T) {
// 	if testing.Short() {
// 		t.SkipNow()
// 	}

// 	if err := godotenv.Load(); err != nil {
// 		log.Fatal(err)
// 	}

// 	text := "Punchinello wanted Payne? He'd see the pain."
// 	entities := []string{"Punchinello", "Payne"}

// 	dps := nlp.NewNLPPoSDetermer(tokenize.ANY)

// 	got, err := Do(nlp, dps, text, entities)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	want := map[tokenize.Token]float64{
// 		{PoS: tokenize.VERB, Token: "wanted"}: 1,
// 		{PoS: tokenize.PUNCT, Token: "?"}:     2,
// 		{PoS: tokenize.PRON, Token: "He"}:     3,
// 		{PoS: tokenize.VERB, Token: "'d"}:     4,
// 		{PoS: tokenize.VERB, Token: "see"}:    5,
// 		{PoS: tokenize.DET, Token: "the"}:     6,
// 		{PoS: tokenize.NOUN, Token: "pain"}:   7,
// 		{PoS: tokenize.PUNCT, Token: "."}:     8,
// 	}
// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("Assoc() = %v, want %v", got, want)
// 	}
// }

// func TestAssocIntegrationMultiWordEntities(t *testing.T) {
// 	if testing.Short() {
// 		t.SkipNow()
// 	}

// 	if err := godotenv.Load(); err != nil {
// 		log.Fatal(err)
// 	}

// 	credentialsFile = os.Getenv("GOOGLE_NLP_SERVICE_ACCOUNT_FILE_LOCATION")

// 	text := "Max Payne, this is Deputy Chief Jim Bravura from the NYPD."
// 	entities := []string{"Max Payne", "Jim Bravura"}

// 	nlp := NewNLP("en")
// 	dps := tokenize.NewPoSDetermer(tokenize.ANY)

// 	got, err := Do(nlp, dps, text, entities)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	want := map[tokenize.Token]float64{
// 		{PoS: tokenize.PUNCT, Token: ","}:     3,
// 		{PoS: tokenize.DET, Token: "this"}:    3,
// 		{PoS: tokenize.VERB, Token: "is"}:     3,
// 		{PoS: tokenize.NOUN, Token: "Deputy"}: 3,
// 		{PoS: tokenize.NOUN, Token: "Chief"}:  3,
// 		{PoS: tokenize.ADP, Token: "from"}:    4,
// 		{PoS: tokenize.DET, Token: "the"}:     5,
// 		{PoS: tokenize.NOUN, Token: "NYPD"}:   6,
// 		{PoS: tokenize.PUNCT, Token: "."}:     7,
// 	}
// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("Assoc() = %v, want %v", got, want)
// 	}
// }

// func TestAssocIntegrationDefinedPartOfSpeech(t *testing.T) {
// 	if testing.Short() {
// 		t.SkipNow()
// 	}

// 	if err := godotenv.Load(); err != nil {
// 		log.Fatal(err)
// 	}

// 	credentialsFile = os.Getenv("GOOGLE_NLP_SERVICE_ACCOUNT_FILE_LOCATION")

// 	text := `"The things that I want", by Max Payne.`
// 	entities := []string{"Max Payne"}

// 	nlp := NewNLP("en")
// 	dps := tokenize.NewPoSDetermer(tokenize.DET | tokenize.VERB | tokenize.PUNCT)

// 	got, err := Do(nlp, dps, text, entities)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	want := map[tokenize.Token]float64{
// 		{PoS: tokenize.PUNCT, Token: `"`}:   4,
// 		{PoS: tokenize.DET, Token: "The"}:   5,
// 		{PoS: tokenize.DET, Token: "that"}:  4,
// 		{PoS: tokenize.VERB, Token: "want"}: 3,
// 		{PoS: tokenize.PUNCT, Token: ","}:   1,
// 		{PoS: tokenize.PUNCT, Token: "."}:   1,
// 	}
// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("Assoc() = %v, want %v", got, want)
// 	}
// }

// // Create a custom NLP instance
// type nlpTest struct{}

// // Second iteration is always for entites
// var iterations int

// func (n *nlpTest) Tokenize(text string) ([]tokenize.Token, error) {
// 	if iterations == 0 {
// 		iterations++

// 		return []tokenize.Token{
// 			{
// 				Token: "Punchinello",
// 				PoS:   tokenize.NOUN,
// 			},
// 			{
// 				Token: "was",
// 				PoS:   tokenize.VERB,
// 			},
// 			{
// 				Token: "burning",
// 				PoS:   tokenize.VERB,
// 			},
// 			{
// 				Token: "to",
// 				PoS:   tokenize.PRT,
// 			},
// 			{
// 				Token: "get",
// 				PoS:   tokenize.VERB,
// 			},
// 			{
// 				Token: "me",
// 				PoS:   tokenize.PRON,
// 			},
// 		}, nil
// 	}

// 	return []tokenize.Token{
// 		{
// 			Token: "Punchinello",
// 			PoS:   tokenize.NOUN,
// 		},
// 	}, nil

// }

// func TestAssocIntegrationSingleWordEntitiesShort(t *testing.T) {
// 	text := "Punchinello was burning to get me"
// 	entities := []string{"Punchinello"}

// 	dps := tokenize.NewPoSDetermer(tokenize.ANY)

// 	got, err := Do(&nlpTest{}, dps, text, entities)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	want := map[tokenize.Token]float64{
// 		{PoS: tokenize.VERB, Token: "was"}:     1,
// 		{PoS: tokenize.VERB, Token: "burning"}: 2,
// 		{PoS: tokenize.PRT, Token: "to"}:       3,
// 		{PoS: tokenize.VERB, Token: "get"}:     4,
// 		{PoS: tokenize.PRON, Token: "me"}:      5,
// 	}
// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("Assoc() = %v, want %v", got, want)
// 	}
// }

// func BenchmarkAssoc(b *testing.B) {
// 	text := "Punchinello was burning to get me"
// 	entities := []string{"Punchinello"}

// 	dps := tokenize.NewPoSDetermer(tokenize.ANY)

// 	for n := 0; n < b.N; n++ {
// 		Do(&nlpTest{}, dps, text, entities)
// 	}
// }
