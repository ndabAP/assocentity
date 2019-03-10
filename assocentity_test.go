package assocentity

import (
	"reflect"
	"testing"
)

func TestRomance_SingleWord(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		entity string
		want   map[string]float64
	}{
		{"empty", "Hello world", "Bye", map[string]float64{}},
		{"subset", "Hello world", "Helloworld", map[string]float64{}},
		{"subset", "Hello world", "Helloworld", map[string]float64{}},
		{"simple", "Hello, my name is Max Payne.", "Max", map[string]float64{
			"Hello": 4,
			"my":    3,
			"name":  2,
			"is":    1,
			"Payne": 1,
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Romance(test.text, test.entity)
			if !reflect.DeepEqual(actual, test.want) {
				t.Errorf("Romance(%v): expected %v, actual %v,", test.name, test.want, actual)
			}
		})
	}
}

func TestRomance_MultiWord(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		entity string
		want   map[string]float64
	}{
		{"empty", "Hello world", "Bye", map[string]float64{}},
		{"subset", "Hello world", "Helloworld", map[string]float64{}},
		{"start", "Shang Tsung is my name.", "Shang Tsung", map[string]float64{
			"is":   1,
			"my":   2,
			"name": 3,
		}},
		{"simple", "Hello, my name is Max Payne.", "Max Payne", map[string]float64{
			"Hello": 4,
			"my":    3,
			"name":  2,
			"is":    1,
		}},
		{
			"inline",
			`If you smell, what Dwayne "The Rock" Johnson is cooking?`,
			`Dwayne "The Rock" Johnson`,
			map[string]float64{
				"If":      4,
				"you":     3,
				"smell":   2,
				"what":    1,
				"is":      1,
				"cooking": 2,
			},
		},
		{
			"inline multi",
			`Shao Kahn is the embodiment of evil. Shao Kahn is easily recognizable by his intimidating stature.`,
			"Shao Kahn",
			map[string]float64{
				"is":           3.75,
				"the":          3,
				"embodiment":   3,
				"of":           3,
				"evil":         3,
				"easily":       5.5,
				"recognizable": 6.5,
				"by":           7.5,
				"his":          8.5,
				"intimidating": 9.5,
				"stature":      10.5,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Romance(test.text, test.entity)
			if !reflect.DeepEqual(actual, test.want) {
				t.Errorf("Romance(%v): expected %v, actual %v,", test.name, test.want, actual)
			}
		})
	}
}

func Test_batch(t *testing.T) {
	type args struct {
		data []int
		size int
	}

	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{"data 4, size 1", args{[]int{1, 2, 3, 4}, 1}, [][]int{{1}, {2}, {3}, {4}}},
		{"data 4, size 2", args{[]int{1, 2, 3, 4}, 2}, [][]int{{1, 2}, {3, 4}}},
		{"data 4, size 3", args{[]int{1, 2, 3, 4}, 3}, [][]int{{1, 2, 3}, {4}}},
		{"data 2, size 3", args{[]int{1, 2}, 3}, [][]int{{1, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := batch(tt.args.data, tt.args.size); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("batch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isSliceSubset(t *testing.T) {
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
			"5 hits",
			args{data: []string{"H", "e", "l", "l", "o"}, subset: []string{"H", "e", "l", "l", "o"}, index: 0},
			true,
		},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSliceSubset(tt.args.data, tt.args.subset, tt.args.index); got != tt.want {
				t.Errorf("isSliceSubset() = %v, want %v", got, tt.want)
			}
		})
	}
}
