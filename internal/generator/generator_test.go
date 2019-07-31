package generator

import (
	"testing"
)

func TestGenerator_Next(t *testing.T) {
	type fields struct {
		slice []string
		pos   int
		el    string
		init  bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "element left",
			fields: fields{
				slice: []string{"Vladimir Lem", "Vincent Gognitti", "Jack Lupino"},
				pos:   0,
				el:    "Vladimir Lem",
				init:  true,
			},
			want: true,
		},
		{
			name: "no element left",
			fields: fields{
				slice: []string{"Vladimir Lem"},
				pos:   0,
				el:    "Vladimir Lem",
				init:  true,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				slice: tt.fields.slice,
				pos:   tt.fields.pos,
				el:    tt.fields.el,
				init:  tt.fields.init,
			}
			if got := g.Next(); got != tt.want {
				t.Errorf("Generator.Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_Prev(t *testing.T) {
	type fields struct {
		slice []string
		pos   int
		el    string
		init  bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "element left",
			fields: fields{
				slice: []string{"Vladimir Lem", "Vincent Gognitti", "Jack Lupino"},
				pos:   2,
				el:    "Jack Lupino",
				init:  true,
			},
			want: true,
		},
		{
			name: "no element left",
			fields: fields{
				slice: []string{"Vladimir Lem"},
				pos:   0,
				el:    "Vladimir Lem",
				init:  true,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				slice: tt.fields.slice,
				pos:   tt.fields.pos,
				el:    tt.fields.el,
				init:  tt.fields.init,
			}
			if got := g.Prev(); got != tt.want {
				t.Errorf("Generator.Prev() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_SetPos(t *testing.T) {
	type fields struct {
		slice []string
		pos   int
		el    string
		init  bool
	}
	type args struct {
		pos int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "element left",
			fields: fields{
				slice: []string{"Vladimir Lem", "Vincent Gognitti", "Jack Lupino"},
				pos:   0,
				el:    "Vladimir Lem",
				init:  true,
			},
			args: args{
				pos: 1,
			},
			want: true,
		},
		{
			name: "no element left",
			fields: fields{
				slice: []string{"Vladimir Lem"},
				pos:   0,
				el:    "Vladimir Lem",
				init:  true,
			},
			args: args{
				pos: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				slice: tt.fields.slice,
				pos:   tt.fields.pos,
				el:    tt.fields.el,
				init:  tt.fields.init,
			}
			if got := g.SetPos(tt.args.pos); got != tt.want {
				t.Errorf("Generator.SetPos() = %v, want %v", got, tt.want)
			}
		})
	}
}
