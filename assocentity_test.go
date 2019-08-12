package assocentity

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	"github.com/ndabAP/assocentity/v5/tokenize"
)

var credentialsFile string

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	credentialsFile = os.Getenv("GOOGLE_NLP_SERVICE_ACCOUNT_FILE_LOCATION")
}

func TestAssocIntegrationSingleWordEntities(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	text := "Punchinello wanted Payne? He'd see the pain."
	entities := []string{"Punchinello", "Payne"}

	nlp, err := tokenize.NewNLP(credentialsFile, text, entities)
	if err != nil {
		log.Fatal(err)
	}

	dps := tokenize.NewPoSDetermer(tokenize.ANY)
	dj := tokenize.NewJoin(tokenize.Whitespace)

	got, err := Do(nlp, dps, dj, entities)
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

func TestAssocIntegrationMultiWordEntities(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	text := "Max Payne, this is Deputy Chief Jim Bravura from the NYPD."
	entities := []string{"Max Payne", "Jim Bravura"}

	nlp, err := tokenize.NewNLP(credentialsFile, text, entities)
	if err != nil {
		log.Fatal(err)
	}

	dps := tokenize.NewPoSDetermer(tokenize.ANY)
	dj := tokenize.NewJoin(tokenize.Whitespace)

	got, err := Do(nlp, dps, dj, entities)
	if err != nil {
		log.Fatal(err)
	}

	want := map[string]float64{
		",":      3,
		"this":   3,
		"Deputy": 3,
		"Chief":  3,
		"is":     3,
		"from":   4,
		"the":    5,
		"NYPD":   6,
		".":      7,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Assoc() = %v, want %v", got, want)
	}
}

func TestAssocIntegrationDefinedPartOfSpeech(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	text := `"The things that I want", by Max Payne.`
	entities := []string{"Max Payne"}

	nlp, err := tokenize.NewNLP(credentialsFile, text, entities)
	if err != nil {
		log.Fatal(err)
	}

	dps := tokenize.NewPoSDetermer(tokenize.DET | tokenize.VERB | tokenize.PUNCT)
	dj := tokenize.NewJoin(tokenize.Whitespace)

	got, err := Do(nlp, dps, dj, entities)
	if err != nil {
		log.Fatal(err)
	}

	want := map[string]float64{
		`"`:    4,
		"The":  5,
		"that": 4,
		"want": 3,
		",":    1,
		".":    1,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Assoc() = %v, want %v", got, want)
	}
}
