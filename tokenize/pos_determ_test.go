package tokenize

import (
	"reflect"
	"testing"
)

type tokenizerTest struct{}

func TestPoSDeterm_Determ(t *testing.T) {
	type fields struct {
		poS int
	}
	type args struct {
		textTokens   []Token
		entityTokens [][]Token
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
				textTokens: []Token{
					{PoS: NOUN, Token: "Cold"},
					{PoS: ADP, Token: "as"},
					{PoS: DET, Token: "a"},
					{PoS: NOUN, Token: "gun"},
				},
				entityTokens: [][]Token{
					{
						{
							Token: "Max",
							PoS:   NOUN,
						},
						{
							Token: "Payne",
							PoS:   NOUN,
						},
					},
				},
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
				textTokens: []Token{
					{PoS: NOUN, Token: "Cold"},
					{PoS: ADP, Token: "as"},
					{PoS: DET, Token: "a"},
					{PoS: NOUN, Token: "gun"},
				},
				entityTokens: [][]Token{
					{
						{
							Token: "Max",
							PoS:   NOUN,
						},
						{
							Token: "Payne",
							PoS:   NOUN,
						},
					},
				},
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
			got, err := dps.Determ(tt.args.textTokens, tt.args.entityTokens)
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
