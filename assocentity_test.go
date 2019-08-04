package assocentity

import (
	"log"
	"reflect"
	"testing"

	"github.com/ndabAP/assocentity/v3/tokenize"
)

type tokenizer struct{}

func (t *tokenizer) TokenizeText() ([]string, error) {
	return texts[currPos], nil
}

func (t *tokenizer) TokenizeEntities() ([][]string, error) {
	return entities[currPos], nil
}

type joiner struct {
	tokens []string
}

func (j *joiner) Join(t tokenize.Tokenizer) error {
	j.tokens = joins[currPos]

	return nil
}

func (j *joiner) Tokens() []string {
	return j.tokens
}

var texts = [][]string{
	{"Vinnie", "Gognitti", "Just", "the", "man", "I", "'ve", "been", "killing", "to", "see", "Gognitti", "bailed"},
	{"You", "can", "'t'", "win", "this", "one", "Max"},
	{"Alex", "Alex"},
	{"I", "'m", "Frankie", "The", "Bat", "Niagara"},
	{"Where", "'s", "Lupino", "Bad", "start", "Vinnie"},
}
var entities = [][][]string{
	{{"Vinnie", "Gognitti"}, {"Gognitti"}},
	{{"Max"}},
	{{"Alex"}},
	{{"Frankie", "The", "Bat", "Niagara"}},
	{{"Lupino"}, {"Vinnie"}},
}
var joins = [][]string{
	{"Vinnie Gognitti", "Just", "the", "man", "I", "'ve", "been", "killing", "to", "see", "Gognitti", "bailed"},
	{"You", "can", "'t", "win", "this", "one", "Max"},
	{"Alex", "Alex"},
	{"I", "'m", "Frankie The Bat Niagara"},
	{"Where", "'s", "Lupino", "Bad", "start", "Vinnie"},
}

var currPos int

func TestAssoc(t *testing.T) {
	type args struct {
		j         tokenize.Joiner
		tokenizer tokenize.Tokenizer
		entities  []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]float64
		wantErr bool
	}{
		{
			name: "two entities",
			args: args{
				j:         new(joiner),
				tokenizer: new(tokenizer),
				entities:  []string{"Vinnie Gognitti", "Gognitti"},
			},
			want: map[string]float64{
				"Just":    5,
				"bailed":  6,
				"to":      5,
				"man":     5,
				"'ve":     5,
				"see":     5,
				"the":     5,
				"I":       5,
				"been":    5,
				"killing": 5,
			},
			wantErr: false,
		},
		{
			name: "one entity (multiple words)",
			args: args{
				j:         new(joiner),
				tokenizer: new(tokenizer),
				entities:  []string{"Max"},
			},
			want: map[string]float64{
				"You":  6,
				"can":  5,
				"'t":   4,
				"win":  3,
				"this": 2,
				"one":  1,
			},
			wantErr: false,
		},
		{
			name: "one entity (only entities)",
			args: args{
				j:         new(joiner),
				tokenizer: new(tokenizer),
				entities:  []string{"Alex"},
			},
			want:    map[string]float64{},
			wantErr: false,
		},
		{
			name: "one entity (multiple tokens)",
			args: args{
				j:         new(joiner),
				tokenizer: new(tokenizer),
				entities:  []string{"Frankie The Bat Niagara"},
			},
			want: map[string]float64{
				"I":  2,
				"'m": 1,
			},
			wantErr: false,
		},
		{
			name: "two entities (in-between)",
			args: args{
				j:         new(joiner),
				tokenizer: new(tokenizer),
				entities:  []string{"Lupino", "Vinnie"},
			},
			want: map[string]float64{
				"Where": 3.5,
				"'s":    2.5,
				"Bad":   1.5,
				"start": 1.5,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Assoc(tt.args.j, tt.args.tokenizer, tt.args.entities)
			if (err != nil) != tt.wantErr {
				t.Errorf("Assoc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Assoc() = %v, want %v", got, tt.want)
			}
		})

		currPos++
	}
}

func TestAssocIntegration1(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	const (
		credentialsFile = "configs/google_nlp_service_account.json"
		sep             = " "
	)

	text := "Punchinello wanted Payne? He'd see the pain."
	entities := []string{"Punchinello", "Payne"}

	nlp, err := tokenize.NewNLP(credentialsFile, text, entities)
	if err != nil {
		log.Fatal(err)
	}

	dj := tokenize.NewDefaultJoin(sep)
	if err = dj.Join(nlp); err != nil {
		log.Fatal(err)
	}

	got, err := Assoc(dj, nlp, entities)
	if err != nil {
		log.Fatal(err)
	}

	wanted := map[string]float64{
		"wanted": 1,
		"?":      2,
		"He":     3,
		"'d":     4,
		"see":    5,
		"the":    6,
		"pain":   7,
		".":      8,
	}
	if !reflect.DeepEqual(got, wanted) {
		t.Errorf("Assoc() = %v, want %v", got, wanted)
	}
}

func TestAssocIntegration2(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	const (
		credentialsFile = "configs/google_nlp_service_account.json"
		sep             = " "
	)

	text := "Max Payne, this is Deputy Chief Jim Bravura from the NYPD."
	entities := []string{"Max Payne", "Jim Bravura"}

	nlp, err := tokenize.NewNLP(credentialsFile, text, entities)
	if err != nil {
		log.Fatal(err)
	}

	dj := tokenize.NewDefaultJoin(sep)
	if err = dj.Join(nlp); err != nil {
		log.Fatal(err)
	}

	got, err := Assoc(dj, nlp, entities)
	if err != nil {
		log.Fatal(err)
	}

	wanted := map[string]float64{
		",":      3,
		".":      7,
		"Chief":  3,
		"Deputy": 3,
		"NYPD":   6,
		"from":   4,
		"is":     3,
		"the":    5,
		"this":   3,
	}
	if !reflect.DeepEqual(got, wanted) {
		t.Errorf("Assoc() = %v, want %v", got, wanted)
	}
}
