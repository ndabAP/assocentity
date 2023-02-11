package assocentity

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/ndabAP/assocentity/v12/tokenize"
)

// whiteSpaceTokenizer tokenizes a text by empty space and assigns unknown
// pos
type whiteSpaceTokenizer int

func (t whiteSpaceTokenizer) Tokenize(ctx context.Context, text string) ([]tokenize.Token, error) {
	spl := strings.Split(text, " ")
	tokens := make([]tokenize.Token, 0)
	for _, s := range spl {
		tokens = append(tokens, tokenize.Token{
			PoS:  tokenize.UNKN,
			Text: s,
		})
	}

	return tokens, nil
}

func TestMean(t *testing.T) {
	type args struct {
		ctx       context.Context
		tokenizer tokenize.Tokenizer
		poS       tokenize.PoS
		texts     []string
		entities  []string
	}
	tests := []struct {
		args args
		want map[[2]string]float64
	}{
		{
			args: args{
				ctx:       context.Background(),
				tokenizer: new(whiteSpaceTokenizer),
				poS:       tokenize.ANY,
				texts: []string{
					"AA B $ CCC ++",
					"$ E ++ AA $ B",
				},
				entities: []string{"$", "++"},
			},
			want: map[[2]string]float64{
				{"AA", tokenize.PoSMapStr[tokenize.UNKN]}:  2.2,
				{"B", tokenize.PoSMapStr[tokenize.UNKN]}:   2.6,
				{"CCC", tokenize.PoSMapStr[tokenize.UNKN]}: 1,
				{"E", tokenize.PoSMapStr[tokenize.UNKN]}:   1.6666666666666667,
			},
		},
		{
			args: args{
				ctx:       context.Background(),
				tokenizer: new(whiteSpaceTokenizer),
				poS:       tokenize.ANY,
				texts: []string{
					"",
					"",
				},
				entities: []string{},
			},
			want: map[[2]string]float64{},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			s := Source{
				Entities: tt.args.entities,
				Texts:    tt.args.texts,
			}
			dists, err := Distances(
				tt.args.ctx,
				tt.args.tokenizer,
				tt.args.poS,
				s,
			)
			if err != nil {
				t.Error(err)
			}

			got := Mean(dists)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mean() = %v, want %v", got, tt.want)
			}
		})
	}
}

// concreteTokenizer is a tokenizer with a fixed set of tokens
type concreteTokenizer int

func (t concreteTokenizer) Tokenize(ctx context.Context, text string) ([]tokenize.Token, error) {
	spl := strings.Split(text, " ")
	tokens := make([]tokenize.Token, 0)
	for _, s := range spl {
		var poS tokenize.PoS
		switch s {
		case "English":
			poS = tokenize.NOUN

		case ".":
			poS = tokenize.PUNCT

		case "run":
			poS = tokenize.VERB

		default:
			continue
		}

		tokens = append(tokens, tokenize.Token{
			PoS:  poS,
			Text: s,
		})
	}

	return tokens, nil
}

func Test_dist(t *testing.T) {
	type args struct {
		ctx       context.Context
		tokenizer tokenize.Tokenizer
		poS       tokenize.PoS
		texts     []string
		entities  []string
	}
	tests := []struct {
		args    args
		want    map[[2]string][]float64
		wantErr bool
	}{
		{
			args: args{
				ctx:       context.Background(),
				tokenizer: new(concreteTokenizer),
				poS:       tokenize.NOUN | tokenize.PUNCT | tokenize.VERB,
				texts:     []string{"English x . x xx run"},
				entities:  []string{"run"},
			},
			want: map[[2]string][]float64{
				{"English", tokenize.PoSMapStr[tokenize.NOUN]}: {2},
				{".", tokenize.PoSMapStr[tokenize.PUNCT]}:      {1},
			},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			s := Source{
				Entities: tt.args.entities,
				Texts:    tt.args.texts,
			}
			got, err := Distances(
				tt.args.ctx,
				tt.args.tokenizer,
				tt.args.poS,
				s,
			)
			if err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dist() = %v, want %v", got, tt.want)
			}
		})
	}
}
