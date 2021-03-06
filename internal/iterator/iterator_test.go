package iterator

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		slice []Element
	}
	tests := []struct {
		name string
		args args
		want *Iterator
	}{
		{
			name: "new",
			args: args{
				slice: []Element{"Gognitti", "bailed"},
			},
			want: &Iterator{
				slice: []Element{"Gognitti", "bailed"},
				pos:   0,
				el:    "Gognitti",
				init:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_Next(t *testing.T) {
	tests := []struct {
		name string
		g    *Iterator
		want bool
	}{
		{
			name: "element left",
			g: &Iterator{
				slice: []Element{"No", "Payne", "No", "Gain"},
				pos:   0,
				el:    "No",
				init:  true,
			},
			want: true,
		},
		{
			name: "no element left",
			g: &Iterator{
				slice: []Element{"No", "Payne", "No", "Gain"},
				pos:   3,
				el:    "Gain",
				init:  false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Next(); got != tt.want {
				t.Errorf("Iterator.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_Prev(t *testing.T) {
	tests := []struct {
		name string
		g    *Iterator
		want bool
	}{
		{
			name: "next",
			g: &Iterator{
				slice: []Element{"No", "Payne", "No", "Gain"},
				pos:   0,
				el:    "No",
				init:  true,
			},
			want: true,
		},
		{
			name: "no element left",
			g: &Iterator{
				slice: []Element{"No", "Payne", "No", "Gain"},
				pos:   0,
				el:    "No",
				init:  false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Prev(); got != tt.want {
				t.Errorf("Iterator.Prev() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_CurrPos(t *testing.T) {
	tests := []struct {
		name string
		g    *Iterator
		next bool
		want int
	}{
		{
			name: "current position",
			g:    New([]Element{"You", "play", "you", "pay"}),
			next: false,
			want: 0,
		},
		{
			name: "next position",
			g:    New([]Element{"You", "play", "you", "pay"}),
			next: true,
			want: 0,
		},
	}
	for _, tt := range tests {
		if tt.next {
			tt.g.Next()
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.CurrPos(); got != tt.want {
				t.Errorf("Iterator.CurrPos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_CurrElem(t *testing.T) {
	tests := []struct {
		name string
		g    *Iterator
		want Element
	}{
		{
			name: "current element",
			g: &Iterator{
				slice: []Element{"Relax", "Max"},
				pos:   0,
				el:    "Relax",
				init:  true,
			},
			want: "Relax",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.CurrElem(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Iterator.CurrElem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_Len(t *testing.T) {
	tests := []struct {
		name string
		g    *Iterator
		want int
	}{
		{
			name: "current element",
			g: &Iterator{
				slice: []Element{"Cold", "as", "a", "gun"},
				pos:   0,
				el:    "Cold",
				init:  true,
			},
			want: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.Len(); got != tt.want {
				t.Errorf("Iterator.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIterator_SetPos(t *testing.T) {
	type args struct {
		pos int
	}
	tests := []struct {
		name string
		g    *Iterator
		args args
		want bool
	}{
		{
			name: "available position",
			g: &Iterator{
				slice: []Element{"With", "pleasure", "boss"},
				pos:   0,
				el:    "Cold",
				init:  true,
			},
			args: args{
				pos: 1,
			},
			want: true,
		},
		{
			name: "unavailable position",
			g: &Iterator{
				slice: []Element{"With", "pleasure", "boss"},
				pos:   2,
				el:    "Cold",
				init:  false,
			},
			args: args{
				pos: 3,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.g.SetPos(tt.args.pos); got != tt.want {
				t.Errorf("Iterator.SetPos() = %v, want %v", got, tt.want)
			}
		})
	}
}
