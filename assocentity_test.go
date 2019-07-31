package assocentity

import (
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

type joiner struct{}

func (j *joiner) Join(t tokenize.Tokenizer) ([]string, error) {
	return joins[currPos], nil
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
