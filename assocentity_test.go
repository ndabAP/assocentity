package assocentity

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/ndabAP/assocentity/v11/tokenize"
)

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

func TestMeanN(t *testing.T) {
	type args struct {
		ctx       context.Context
		tokenizer tokenize.Tokenizer
		poS       tokenize.PoS
		texts     []string
		entities  []string
	}
	tests := []struct {
		args    args
		want    map[tokenize.Token]float64
		wantErr bool
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
			want: map[tokenize.Token]float64{
				{
					PoS:  tokenize.UNKN,
					Text: "AA",
				}: 2.2,
				{
					PoS:  tokenize.UNKN,
					Text: "B",
				}: 2.6,
				{
					PoS:  tokenize.UNKN,
					Text: "CCC",
				}: 1,
				{
					PoS:  tokenize.UNKN,
					Text: "E",
				}: 1.6666666666666667,
			},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, err := MeanN(tt.args.ctx, tt.args.tokenizer, tt.args.poS, tt.args.texts, tt.args.entities)
			if (err != nil) != tt.wantErr {
				t.Errorf("MeanN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MeanN() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		text      string
		entities  []string
	}
	tests := []struct {
		args    args
		want    map[tokenize.Token][]float64
		wantErr bool
	}{
		{
			args: args{
				ctx:       context.Background(),
				tokenizer: new(concreteTokenizer),
				poS:       tokenize.NOUN | tokenize.PUNCT | tokenize.VERB,
				text:      "English x . x xx run",
				entities:  []string{"run"},
			},
			want: map[tokenize.Token][]float64{
				{
					PoS:  tokenize.NOUN,
					Text: "English",
				}: {2},
				{
					PoS:  tokenize.PUNCT,
					Text: ".",
				}: {1},
			},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got, err := dist(tt.args.ctx, tt.args.tokenizer, tt.args.poS, tt.args.text, tt.args.entities)
			if (err != nil) != tt.wantErr {
				t.Errorf("dist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dist() = %v, want %v", got, tt.want)
			}
		})
	}
}
