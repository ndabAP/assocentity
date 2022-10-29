package comp_test

import (
	"reflect"
	"testing"

	"github.com/ndabAP/assocentity/v9/internal/comp"
	"github.com/ndabAP/assocentity/v9/internal/iterator"
	"github.com/ndabAP/assocentity/v9/tokenize"
)

func TestTextWithEntity(t *testing.T) {
	type args struct {
		textIter         *iterator.Iterator[tokenize.Token]
		entityTokensIter *iterator.Iterator[[]tokenize.Token]
		dir              comp.Direction
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 []tokenize.Token
	}{
		// {
		// 	name: "no entity",
		// 	args: args{
		// 		textIter: iterator.New([]tokenize.Token{
		// 			{
		// 				PoS:  tokenize.ADP,
		// 				Text: "Without",
		// 			},
		// 			{
		// 				PoS:  tokenize.NOUN,
		// 				Text: "Mona",
		// 			},
		// 			{
		// 				PoS:  tokenize.PRT,
		// 				Text: "'s'",
		// 			},
		// 			{
		// 				PoS:  tokenize.NOUN,
		// 				Text: "help",
		// 			},
		// 			{
		// 				PoS:  tokenize.PUNCT,
		// 				Text: ",",
		// 			},
		// 			{
		// 				PoS:  tokenize.PRON,
		// 				Text: "I",
		// 			},
		// 			{
		// 				PoS:  tokenize.VERB,
		// 				Text: "'d'",
		// 			},
		// 			{
		// 				PoS:  tokenize.VERB,
		// 				Text: "be",
		// 			},
		// 			{
		// 				PoS:  tokenize.DET,
		// 				Text: "a",
		// 			},
		// 			{
		// 				PoS:  tokenize.ADJ,
		// 				Text: "dead",
		// 			},
		// 			{
		// 				PoS:  tokenize.NOUN,
		// 				Text: "man",
		// 			},
		// 		}),
		// 		entityTokensIter: iterator.New([][]tokenize.Token{
		// 			{
		// 				{
		// 					PoS:  tokenize.NOUN,
		// 					Text: "Alex",
		// 				},
		// 			},
		// 		}),
		// 		dir: comp.DirPos,
		// 	},
		// 	want:  false,
		// 	want1: make([]tokenize.Token, 0),
		// },
		{
			name: "entity",
			args: args{
				textIter: iterator.New([]tokenize.Token{
					{
						PoS:  tokenize.ADP,
						Text: "Without",
					},
					{
						PoS:  tokenize.NOUN,
						Text: "Mona",
					},
					{
						PoS:  tokenize.PRT,
						Text: "'s'",
					},
					{
						PoS:  tokenize.NOUN,
						Text: "help",
					},
					{
						PoS:  tokenize.PUNCT,
						Text: ",",
					},
					{
						PoS:  tokenize.PRON,
						Text: "I",
					},
					{
						PoS:  tokenize.VERB,
						Text: "'d'",
					},
					{
						PoS:  tokenize.VERB,
						Text: "be",
					},
					{
						PoS:  tokenize.DET,
						Text: "a",
					},
					{
						PoS:  tokenize.ADJ,
						Text: "dead",
					},
					{
						PoS:  tokenize.NOUN,
						Text: "man",
					},
				}).SetPos(1),
				entityTokensIter: iterator.New([][]tokenize.Token{
					{
						tokenize.Token{
							PoS:  tokenize.NOUN,
							Text: "Mona",
						},
					},
				}),
				dir: comp.DirPos,
			},
			want: true,
			want1: []tokenize.Token{
				{
					PoS:  tokenize.NOUN,
					Text: "Mona",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := comp.TextWithEntities(tt.args.textIter, tt.args.entityTokensIter, tt.args.dir)
			if got != tt.want {
				t.Errorf("TextWithEntity() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("TextWithEntity() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
