package assocentity

import "testing"

func TestMake(t *testing.T) {
	type args struct {
		text      string
		entities  []string
		tokenizer func(string) ([]string, error)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Test",
			args{
				"",
				[]string{},
				nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Make(tt.args.text, tt.args.entities, tt.args.tokenizer)
		})
	}
}
