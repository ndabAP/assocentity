package tokenize

import (
	"reflect"
	"testing"
)

type tokenizerTest struct{}

func (t *tokenizerTest) TokenizeText() ([]Token, error) {
	return []Token{
		{
			Token: "Cold",
			PoS:   NOUN,
		},
		{
			Token: "as",
			PoS:   ADP,
		},
		{
			Token: "a",
			PoS:   DET,
		},
		{
			Token: "gun",
			PoS:   NOUN,
		},
	}, nil
}

func (t *tokenizerTest) TokenizeEntities() ([][]Token, error) {
	return [][]Token{
		{
			Token{
				Token: "Max",
				PoS:   NOUN,
			},
			Token{
				Token: "Payne",
				PoS:   NOUN,
			},
		},
	}, nil
}

func TestPoSDeterm_Determ(t *testing.T) {
	type fields struct {
		poS int
	}
	type args struct {
		tokenizer Tokenizer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Token
		wantErr bool
	}{
		{
			name: "any",
			fields: fields{
				poS: ANY,
			},
			args: args{
				tokenizer: new(tokenizerTest),
			},
			want: []Token{
				{PoS: NOUN, Token: "Cold"},
				{PoS: ADP, Token: "as"},
				{PoS: DET, Token: "a"},
				{PoS: NOUN, Token: "gun"},
			},
			wantErr: false,
		},
		{
			name: "noun",
			fields: fields{
				poS: NOUN,
			},
			args: args{
				tokenizer: new(tokenizerTest),
			},
			want: []Token{
				{PoS: NOUN, Token: "Cold"},
				{PoS: NOUN, Token: "gun"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dps := &PoSDeterm{
				poS: tt.fields.poS,
			}
			got, err := dps.Determ(tt.args.tokenizer)
			if (err != nil) != tt.wantErr {
				t.Errorf("PoSDeterm.Determ() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PoSDeterm.Determ() = %v, want %v", got, tt.want)
			}
		})
	}
}
