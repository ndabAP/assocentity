package tokenize

import (
	"reflect"
	"testing"
)

type tokenizer string

func (t *tokenizer) TokenizeText() ([]string, error) {
	return texts[currPos], nil
}

func (t *tokenizer) TokenizeEntities() ([][]string, error) {
	return entities[currPos], nil
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

var currPos int

func TestDefaultJoin_Join(t *testing.T) {
	type args struct {
		tok Tokenizer
	}
	tests := []struct {
		name    string
		dm      *DefaultJoin
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "multiple entities (two occurrences)",
			dm:   &DefaultJoin{sep: " "},
			args: args{
				new(tokenizer),
			},
			want:    []string{"Vinnie Gognitti", "Just", "the", "man", "I", "'ve", "been", "killing", "to", "see", "Gognitti", "bailed"},
			wantErr: false,
		},
		{
			name: "one entity (one occurrence)",
			dm:   &DefaultJoin{sep: " "},
			args: args{
				new(tokenizer),
			},
			want:    []string{"You", "can'", "'t'", "win", "this", "one", "Max"},
			wantErr: false,
		},
		{
			name: "one entity (one occurrence, twice)",
			dm:   &DefaultJoin{sep: " "},
			args: args{
				new(tokenizer),
			},
			want:    []string{"Alex", "Alex"},
			wantErr: false,
		},
		{
			name: "one entity (multiple occurrences)",
			dm:   &DefaultJoin{sep: " "},
			args: args{
				new(tokenizer),
			},
			want:    []string{"I", "'m", "Frankie The Bat Niagara"},
			wantErr: false,
		},
		{
			name: "multiple entities (one occurrence)",
			dm:   &DefaultJoin{sep: " "},
			args: args{
				new(tokenizer),
			},
			want:    []string{"Where", "'s", "Lupino", "Bad", "start", "Vinnie"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dm := tt.dm
			got, err := dm.Join(tt.args.tok)
			if (err != nil) != tt.wantErr {
				t.Errorf("DefaultMultiplex.Multiplex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultMultiplex.Multiplex() = %v, want %v", got, tt.want)
			}
		})

		currPos++
	}
}
