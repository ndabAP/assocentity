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

var nlpTokenizer nlp.NLPTokenizer

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	credentialsFile := os.Getenv("GOOGLE_NLP_SERVICE_ACCOUNT_FILE_LOCATION")
	nlpTokenizer = nlp.NewNLPTokenizer(credentialsFile, nlp.AutoLang)
}

func TestDoSimple1(t *testing.T) {
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
	if testing.Short() {
		t.SkipNow()
	}

	text := "ee ee aa bb cc dd. b ff, gg, hh, bb, bb. ii!"
	entities := []string{"bb", "b", "ee ee"}

	posDeterm := nlp.NewNLPPoSDetermer(tokenize.ANY)

	got, err := assocentity.Do(context.Background(), nlpTokenizer, posDeterm, text, entities)
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
