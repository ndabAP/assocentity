package nlp

import (
	"reflect"
	"testing"

	"github.com/ndabAP/assocentity/v8/tokenize"
)

func TestNLPPoSDetermer_DetermPoS(t *testing.T) {
	type fields struct {
		poS tokenize.PoS
	}
	type args struct {
		textTokens   []tokenize.Token
		entityTokens [][]tokenize.Token
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []tokenize.Token
	}{
		// {
		// 	name: "any",
		// 	fields: fields{
		// 		poS: tokenize.ANY,
		// 	},
		// 	args: args{
		// 		textTokens: []tokenize.Token{
		// 			{PoS: tokenize.NOUN, Text: "Cold"},
		// 			{PoS: tokenize.ADP, Text: "as"},
		// 			{PoS: tokenize.DET, Text: "a"},
		// 			{PoS: tokenize.NOUN, Text: "gun"},
		// 		},
		// 		entityTokens: [][]tokenize.Token{
		// 			{
		// 				{
		// 					Text: "Max",
		// 					PoS:  tokenize.NOUN,
		// 				},
		// 				{
		// 					Text: "Payne",
		// 					PoS:  tokenize.NOUN,
		// 				},
		// 			},
		// 		},
		// 	},
		// 	want: []tokenize.Token{
		// 		{PoS: tokenize.NOUN, Text: "Cold"},
		// 		{PoS: tokenize.ADP, Text: "as"},
		// 		{PoS: tokenize.DET, Text: "a"},
		// 		{PoS: tokenize.NOUN, Text: "gun"},
		// 	},
		// },
		// {
		// 	name: "noun",
		// 	fields: fields{
		// 		poS: tokenize.NOUN,
		// 	},
		// 	args: args{
		// 		textTokens: []tokenize.Token{
		// 			{PoS: tokenize.NOUN, Text: "Cold"},
		// 			{PoS: tokenize.ADP, Text: "as"},
		// 			{PoS: tokenize.DET, Text: "a"},
		// 			{PoS: tokenize.NOUN, Text: "gun"},
		// 		},
		// 		entityTokens: [][]tokenize.Token{
		// 			{
		// 				{
		// 					Text: "Max",
		// 					PoS:  tokenize.NOUN,
		// 				},
		// 				{
		// 					Text: "Payne",
		// 					PoS:  tokenize.NOUN,
		// 				},
		// 			},
		// 		},
		// 	},
		// 	want: []tokenize.Token{
		// 		{PoS: tokenize.NOUN, Text: "Cold"},
		// 		{PoS: tokenize.NOUN, Text: "gun"},
		// 	},
		// },
		// {
		// 	name: "noun, adposition",
		// 	fields: fields{
		// 		poS: tokenize.NOUN | tokenize.ADP,
		// 	},
		// 	args: args{
		// 		textTokens: []tokenize.Token{
		// 			{PoS: tokenize.NOUN, Text: "Cold"},
		// 			{PoS: tokenize.ADP, Text: "as"},
		// 			{PoS: tokenize.DET, Text: "a"},
		// 			{PoS: tokenize.NOUN, Text: "gun"},
		// 		},
		// 		entityTokens: [][]tokenize.Token{
		// 			{
		// 				{
		// 					Text: "Max",
		// 					PoS:  tokenize.NOUN,
		// 				},
		// 				{
		// 					Text: "Payne",
		// 					PoS:  tokenize.NOUN,
		// 				},
		// 			},
		// 		},
		// 	},
		// 	want: []tokenize.Token{
		// 		{PoS: tokenize.NOUN, Text: "Cold"},
		// 		{PoS: tokenize.ADP, Text: "as"},
		// 		{PoS: tokenize.NOUN, Text: "gun"},
		// 	},
		// },
		{
			name: "skip entity",
			fields: fields{
				poS: tokenize.VERB,
			},
			args: args{
				textTokens: []tokenize.Token{
					{PoS: tokenize.VERB, Text: "Relax"},
					{PoS: tokenize.PUNCT, Text: ","},
					{PoS: tokenize.NOUN, Text: "Max"},
					{PoS: tokenize.PUNCT, Text: "."},
					{PoS: tokenize.PRON, Text: "You"},
					{PoS: tokenize.VERB, Text: "'re"},
					{PoS: tokenize.DET, Text: "a"},
					{PoS: tokenize.ADJ, Text: "nice"},
					{PoS: tokenize.NOUN, Text: "guy"},
					{PoS: tokenize.PUNCT, Text: "."},
				},
				entityTokens: [][]tokenize.Token{
					{
						{
							Text: "Max",
							PoS:  tokenize.NOUN,
						},
						{
							Text: "Payne",
							PoS:  tokenize.NOUN,
						},
					},
					{
						{
							Text: "Max",
							PoS:  tokenize.NOUN,
						},
					},
				},
			},
			want: []tokenize.Token{
				{PoS: tokenize.VERB, Text: "Relax"},
				{PoS: tokenize.NOUN, Text: "Max"},
				{PoS: tokenize.VERB, Text: "'re"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dps := NLPPoSDetermer{
				poS: tt.fields.poS,
			}
			if got := dps.DetermPoS(tt.args.textTokens, tt.args.entityTokens); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NLPPoSDetermer.DetermPoS() = %v, want %v", got, tt.want)
			}
		})
	}
}
