package tokenize

import (
	"reflect"
	"testing"
)

type tokenizer string

func (t *tokenizer) TokenizeText() ([]string, error) {
	return texts[pointer], nil
}

func (t *tokenizer) TokenizeEntities() ([][]string, error) {
	return entities[pointer], nil
}

var texts = [][]string{
	{"Vinnie", "Gognitti", "Just", "the", "man", "I", "'ve", "been", "killing", "to", "see", "Gognitti", "bailed"},
	{"You", "can'", "'t'", "win", "this", "one", "Max"},
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

var pointer int

func TestDefaultMultiplex_Multiplex(t *testing.T) {
	type args struct {
		tok Tokenizer
	}
	tests := []struct {
		name    string
		dm      *DefaultMultiplex
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "multiple entities (two occurrences)",
			dm:   &DefaultMultiplex{},
			args: args{
				new(tokenizer),
			},
			want:    []string{"Vinnie Gognitti", "Just", "the", "man", "I", "'ve", "been", "killing", "to", "see", "Gognitti", "bailed"},
			wantErr: false,
		},
		{
			name: "one entity (one occurrence)",
			dm:   &DefaultMultiplex{},
			args: args{
				new(tokenizer),
			},
			want:    []string{"You", "can'", "'t'", "win", "this", "one", "Max"},
			wantErr: false,
		},
		{
			name: "one entity (one occurrence, twice)",
			dm:   &DefaultMultiplex{},
			args: args{
				new(tokenizer),
			},
			want:    []string{"Alex", "Alex"},
			wantErr: false,
		},
		{
			name: "one entity (multiple occurrences)",
			dm:   &DefaultMultiplex{},
			args: args{
				new(tokenizer),
			},
			want:    []string{"I", "'m", "Frankie The Bat Niagara"},
			wantErr: false,
		},
		{
			name: "multiple entities (one occurrence)",
			dm:   &DefaultMultiplex{},
			args: args{
				new(tokenizer),
			},
			want:    []string{"Where", "'s", "Lupino", "Bad", "start", "Vinnie"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dm := &DefaultMultiplex{}
			got, err := dm.Multiplex(tt.args.tok)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultMultiplex.Multiplex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultMultiplex.Multiplex() = %v, want %v", got, tt.want)
			}
		})

		pointer++
	}
}
