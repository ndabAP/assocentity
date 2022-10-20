package nlp

// import (
// 	"reflect"
// 	"testing"

// 	"github.com/ndabAP/assocentity/v8/tokenize"
// )

// type tokenizerTest struct{}

// func TestPoSDeterm_Determ(t *testing.T) {
// 	type fields struct {
// 		poS tokenize.PoS
// 	}
// 	type args struct {
// 		textTokens   []tokenize.Token
// 		entityTokens [][]tokenize.Token
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    []tokenize.Token
// 		wantErr bool
// 	}{
// 		{
// 			name: "any",
// 			fields: fields{
// 				poS: tokenize.ANY,
// 			},
// 			args: args{
// 				textTokens: []tokenize.Token{
// 					{PoS: tokenize.NOUN, Token: "Cold"},
// 					{PoS: tokenize.ADP, Token: "as"},
// 					{PoS: tokenize.DET, Token: "a"},
// 					{PoS: tokenize.NOUN, Token: "gun"},
// 				},
// 				entityTokens: [][]tokenize.Token{
// 					{
// 						{
// 							Token: "Max",
// 							PoS:   tokenize.NOUN,
// 						},
// 						{
// 							Token: "Payne",
// 							PoS:   tokenize.NOUN,
// 						},
// 					},
// 				},
// 			},
// 			want: []tokenize.Token{
// 				{PoS: tokenize.NOUN, Token: "Cold"},
// 				{PoS: tokenize.ADP, Token: "as"},
// 				{PoS: tokenize.DET, Token: "a"},
// 				{PoS: tokenize.NOUN, Token: "gun"},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "noun",
// 			fields: fields{
// 				poS: tokenize.NOUN,
// 			},
// 			args: args{
// 				textTokens: []tokenize.Token{
// 					{PoS: tokenize.NOUN, Token: "Cold"},
// 					{PoS: tokenize.ADP, Token: "as"},
// 					{PoS: tokenize.DET, Token: "a"},
// 					{PoS: tokenize.NOUN, Token: "gun"},
// 				},
// 				entityTokens: [][]tokenize.Token{
// 					{
// 						{
// 							Token: "Max",
// 							PoS:   tokenize.NOUN,
// 						},
// 						{
// 							Token: "Payne",
// 							PoS:   tokenize.NOUN,
// 						},
// 					},
// 				},
// 			},
// 			want: []tokenize.Token{
// 				{PoS: tokenize.NOUN, Token: "Cold"},
// 				{PoS: tokenize.NOUN, Token: "gun"},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			dps := NLPPoSDetermer{
// 				poS: tt.fields.poS,
// 			}
// 			got, err := dps.DetermPoS(tt.args.textTokens, tt.args.entityTokens)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("PoSDeterm.Determ() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("PoSDeterm.Determ() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
