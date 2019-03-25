// Package assocentity returns the average distance from words to a given entity.
package assocentity

import (
	"reflect"
	"strings"
	"testing"
)

func TestEnglish(t *testing.T) {
	type args struct {
		text      string
		entities  []string
		tokenizer Tokenizer
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]float64
		wantErr bool
	}{
		{
			"No aliases",
			args{
				"The quick brown fox jumps over the lazy dog. It seems that this fox has nothing to loose.",
				[]string{"fox"},
				nil,
			},
			map[string]float64{
				"The":     8.5,
				"quick":   7.5,
				"brown":   6.5,
				"jumps":   5.5,
				"over":    5.5,
				"the":     5.5,
				"lazy":    5.5,
				"dog":     5.5,
				".":       8,
				"It":      5.5,
				"seems":   5.5,
				"that":    5.5,
				"this":    5.5,
				"has":     6.5,
				"nothing": 7.5,
				"to":      8.5,
				"loose":   9.5,
			},
			false,
		},
		{
			"Aliases",
			args{
				"The quick brown fox jumps over the lazy dog. It seems that this fox has nothing to loose.",
				[]string{"brown fox", "fox"},
				nil,
			},
			map[string]float64{
				"seems":   5.5,
				"has":     6.5,
				"It":      5.5,
				".":       8,
				"nothing": 7.5,
				"the":     5.5,
				"to":      8.5,
				"dog":     5.5,
				"The":     8,
				"over":    5.5,
				"jumps":   5.5,
				"that":    5.5,
				"lazy":    5.5,
				"loose":   9.5,
				"this":    5.5,
				"quick":   7,
			},
			false,
		},
		{
			"Error: Entity not found",
			args{
				"The quick brown fox jumps over the lazy dog. It seems that this fox has nothing to loose.",
				[]string{"cat"},
				nil,
			},
			nil,
			true,
		},
		{
			"Custom tokenizer",
			args{
				"Hello, world",
				[]string{"world"},
				func(text string) ([]string, error) {
					return strings.Split(strings.Replace(text, " ", "", -1), ","), nil
				},
			},
			map[string]float64{
				"Hello": 1,
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Make(tt.args.text, tt.args.entities, tt.args.tokenizer)
			if (err != nil) != tt.wantErr {
				t.Errorf("Make() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Make() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSliceSubset(t *testing.T) {
	type args struct {
		data   []string
		subset []string
		index  int
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"0 hits",
			args{data: []string{"H", "e", "l", "l", "o"}, subset: []string{"H"}, index: 1},
			false,
		},
		{
			"1 hit",
			args{data: []string{"H", "e", "l", "l", "o"}, subset: []string{"e"}, index: 1},
			true,
		},
		{
			"5 hits",
			args{data: []string{"H", "e", "l", "l", "o"}, subset: []string{"H", "e", "l", "l", "o"}, index: 0},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSliceSubset(tt.args.data, tt.args.subset, tt.args.index); got != tt.want {
				t.Errorf("isSliceSubset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAverage(t *testing.T) {
	tests := []struct {
		name string
		args []float64
		want float64
	}{
		{"Average of 1, 2, 3", []float64{1, 2, 3}, 2},
		{"Average of 3, 2, 1", []float64{3, 2, 1}, 2},
		{"Average of 1", []float64{1}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := average(tt.args); got != tt.want {
				t.Errorf("average() = %v, want %v", got, tt.want)
			}
		})
	}
}
