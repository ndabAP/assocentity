package assocentity

import (
	"reflect"
	"testing"
)

func Test_tokenize(t *testing.T) {
	tests := []struct {
		name string
		text string
		want []string
	}{
		{"1 word", "Hello", []string{"Hello"}},
		{"2 words", "Hello world", []string{"Hello", "world"}},
		{"2 words, punctuation", "Hello world!", []string{"Hello", "world"}},
		{"1 word, dash", "Hello-world", []string{"Hello-world"}},
		{"2 words, dash", "Your hello-world", []string{"Your", "hello-world"}},
		{"1 word, apostrophe", "I'm", []string{"I'm"}},
		{"2 words, apostrophe", "I'm here", []string{"I'm", "here"}},
		{"3 words", "I'm here today", []string{"I'm", "here", "today"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tokenize(tt.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
