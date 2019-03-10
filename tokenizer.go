package assocentity

import (
	"strings"
	"text/scanner"
	"unicode"
)

// Slices a text by white space
func tokenize(text string) []string {
	var s scanner.Scanner

	s.Init(strings.NewReader(text))
	s.IsIdentRune = func(ch rune, i int) bool {
		return ch == escapedapos || ch == dashchar || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i >= 0
	}

	var words []string
	for token := s.Scan(); token != scanner.EOF; token = s.Scan() {
		isPunct := false
		for _, r := range s.TokenText() {
			if unicode.IsPunct(r) && r != unicodeapostrophe && r != uncodedash {
				isPunct = true
			}
		}

		if !isPunct {
			words = append(words, s.TokenText())
		}
	}

	return words
}
