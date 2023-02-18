package assocentity

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/ndabAP/assocentity/v13/tokenize"
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
			source := NewSource(tt.args.entities, tt.args.texts)
			dists, err := Distances(
				tt.args.ctx,
				tt.args.tokenizer,
				tt.args.poS,
				source,
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

func Test_distances(t *testing.T) {
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
			got, err := distances(
				tt.args.ctx,
				tt.args.tokenizer,
				tt.args.poS,
				tt.args.text,
				tt.args.entities,
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

func TestNormalize(t *testing.T) {
	t.Run("HumandReadableNormalizer", func(t *testing.T) {
		got := map[tokenize.Token][]float64{
			{
				PoS:  tokenize.UNKN,
				Text: "A",
			}: {},
			{
				PoS:  tokenize.UNKN,
				Text: "a",
			}: {},
			{
				PoS:  tokenize.UNKN,
				Text: "b",
			}: {},
			{
				PoS:  tokenize.UNKN,
				Text: "&",
			}: {},
		}
		want := map[tokenize.Token][]float64{
			{
				PoS:  tokenize.UNKN,
				Text: "a",
			}: {},
			{
				PoS:  tokenize.UNKN,
				Text: "b",
			}: {},
			{
				PoS:  tokenize.UNKN,
				Text: "and",
			}: {},
		}
		Normalize(got, HumandReadableNormalizer)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Normalize() = %v, want %v", got, want)
		}
	})
}

func TestThreshold(t *testing.T) {
	type args struct {
		dists     map[tokenize.Token][]float64
		threshold float64
	}
	tests := []struct {
		args args
		want map[tokenize.Token][]float64
	}{
		{
			args: args{
				dists: map[tokenize.Token][]float64{
					{
						PoS:  tokenize.UNKN,
						Text: "A",
					}: {1},
					{
						PoS:  tokenize.UNKN,
						Text: "B",
					}: {1, 1},
					{
						PoS:  tokenize.UNKN,
						Text: "C",
					}: {1, 1, 1},
					{
						PoS:  tokenize.UNKN,
						Text: "D",
					}: {1, 1, 1},
				},
				threshold: 75,
			},
			want: map[tokenize.Token][]float64{
				{
					PoS:  tokenize.UNKN,
					Text: "C",
				}: {1, 1, 1},
				{
					PoS:  tokenize.UNKN,
					Text: "D",
				}: {1, 1, 1},
			},
		},
		{
			args: args{
				dists: map[tokenize.Token][]float64{
					{
						PoS:  tokenize.UNKN,
						Text: "A",
					}: {1},
					{
						PoS:  tokenize.UNKN,
						Text: "B",
					}: {1, 1},
					{
						PoS:  tokenize.UNKN,
						Text: "C",
					}: {1, 1, 1},
					{
						PoS:  tokenize.UNKN,
						Text: "D",
					}: {1, 1, 1, 1},
				},
				threshold: 76,
			},
			want: map[tokenize.Token][]float64{
				{
					PoS:  tokenize.UNKN,
					Text: "D",
				}: {1, 1, 1, 1},
			},
		},
		{
			args: args{
				dists: map[tokenize.Token][]float64{
					{
						PoS:  tokenize.UNKN,
						Text: "A",
					}: {1},
					{
						PoS:  tokenize.UNKN,
						Text: "B",
					}: {1, 1},
					{
						PoS:  tokenize.UNKN,
						Text: "C",
					}: {1, 1, 1},
					{
						PoS:  tokenize.UNKN,
						Text: "D",
					}: {1, 1, 1, 1},
				},
				threshold: 1,
			},
			want: map[tokenize.Token][]float64{
				{
					PoS:  tokenize.UNKN,
					Text: "A",
				}: {1},
				{
					PoS:  tokenize.UNKN,
					Text: "B",
				}: {1, 1},
				{
					PoS:  tokenize.UNKN,
					Text: "C",
				}: {1, 1, 1},
				{
					PoS:  tokenize.UNKN,
					Text: "D",
				}: {1, 1, 1, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			Threshold(tt.args.dists, tt.args.threshold)
			if !reflect.DeepEqual(tt.args.dists, tt.want) {
				t.Errorf("Threshold() = %v, want %v", tt.args.dists, tt.want)
			}
		})
	}
}
