package tokenize

import (
	"reflect"
	"testing"
)

func TestNewJoin(t *testing.T) {
	type args struct {
		sep string
	}
	tests := []struct {
		name string
		args args
		want *Join
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJoin(tt.args.sep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJoin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJoin_Join(t *testing.T) {
	type fields struct {
		sep string
	}
	type args struct {
		dps       PoSDetermer
		tokenizer Tokenizer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dj := &Join{
				sep: tt.fields.sep,
			}
			got, err := dj.Join(tt.args.dps, tt.args.tokenizer)
			if (err != nil) != tt.wantErr {
				t.Errorf("Join.Join() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Join.Join() = %v, want %v", got, tt.want)
			}
		})
	}
}
