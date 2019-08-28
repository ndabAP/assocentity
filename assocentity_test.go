package assocentity

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	"github.com/ndabAP/assocentity/v6/tokenize"
)

var credentialsFile string

func TestAssocIntegrationSingleWordEntities(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	credentialsFile = os.Getenv("GOOGLE_NLP_SERVICE_ACCOUNT_FILE_LOCATION")

	text := "Punchinello wanted Payne? He'd see the pain."
	entities := []string{"Punchinello", "Payne"}

	nlp, err := tokenize.NewNLP(credentialsFile, text, entities)
	if err != nil {
		log.Fatal(err)
	}

	dps := tokenize.NewPoSDetermer(tokenize.ANY)

	got, err := Do(nlp, dps, entities)
	if err != nil {
		log.Fatal(err)
	}

	want := map[tokenize.Token]float64{
		tokenize.Token{PoS: tokenize.VERB, Token: "wanted"}: 1,
		tokenize.Token{PoS: tokenize.PUNCT, Token: "?"}:     2,
		tokenize.Token{PoS: tokenize.PRON, Token: "He"}:     3,
		tokenize.Token{PoS: tokenize.VERB, Token: "'d"}:     4,
		tokenize.Token{PoS: tokenize.VERB, Token: "see"}:    5,
		tokenize.Token{PoS: tokenize.DET, Token: "the"}:     6,
		tokenize.Token{PoS: tokenize.NOUN, Token: "pain"}:   7,
		tokenize.Token{PoS: tokenize.PUNCT, Token: "."}:     8,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Assoc() = %v, want %v", got, want)
	}
}

func TestAssocIntegrationMultiWordEntities(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	credentialsFile = os.Getenv("GOOGLE_NLP_SERVICE_ACCOUNT_FILE_LOCATION")

	text := "Max Payne, this is Deputy Chief Jim Bravura from the NYPD."
	entities := []string{"Max Payne", "Jim Bravura"}

	nlp, err := tokenize.NewNLP(credentialsFile, text, entities)
	if err != nil {
		log.Fatal(err)
	}

	dps := tokenize.NewPoSDetermer(tokenize.ANY)

	got, err := Do(nlp, dps, entities)
	if err != nil {
		log.Fatal(err)
	}

	want := map[tokenize.Token]float64{
		tokenize.Token{PoS: tokenize.PUNCT, Token: ","}:     3,
		tokenize.Token{PoS: tokenize.DET, Token: "this"}:    3,
		tokenize.Token{PoS: tokenize.VERB, Token: "is"}:     3,
		tokenize.Token{PoS: tokenize.NOUN, Token: "Deputy"}: 3,
		tokenize.Token{PoS: tokenize.NOUN, Token: "Chief"}:  3,
		tokenize.Token{PoS: tokenize.ADP, Token: "from"}:    4,
		tokenize.Token{PoS: tokenize.DET, Token: "the"}:     5,
		tokenize.Token{PoS: tokenize.NOUN, Token: "NYPD"}:   6,
		tokenize.Token{PoS: tokenize.PUNCT, Token: "."}:     7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Assoc() = %v, want %v", got, want)
	}
}

func TestAssocIntegrationDefinedPartOfSpeech(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	credentialsFile = os.Getenv("GOOGLE_NLP_SERVICE_ACCOUNT_FILE_LOCATION")

	text := `"The things that I want", by Max Payne.`
	entities := []string{"Max Payne"}

	nlp, err := tokenize.NewNLP(credentialsFile, text, entities)
	if err != nil {
		log.Fatal(err)
	}

	dps := tokenize.NewPoSDetermer(tokenize.DET | tokenize.VERB | tokenize.PUNCT)

	got, err := Do(nlp, dps, entities)
	if err != nil {
		log.Fatal(err)
	}

	want := map[tokenize.Token]float64{
		tokenize.Token{PoS: tokenize.PUNCT, Token: `"`}:   4,
		tokenize.Token{PoS: tokenize.DET, Token: "The"}:   5,
		tokenize.Token{PoS: tokenize.DET, Token: "that"}:  4,
		tokenize.Token{PoS: tokenize.VERB, Token: "want"}: 3,
		tokenize.Token{PoS: tokenize.PUNCT, Token: ","}:   1,
		tokenize.Token{PoS: tokenize.PUNCT, Token: "."}:   1,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Assoc() = %v, want %v", got, want)
	}
}
