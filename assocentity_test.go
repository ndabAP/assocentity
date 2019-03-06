package assocentity

import (
	"reflect"
	"testing"
)

func TestLatin_simpleOneWord(t *testing.T) {
	text := "Hello, my name is Max Payne."
	entity := "Max"

	res := Latin(text, entity)
	m := map[string]float64{
		"Hello": 4,
		"my":    3,
		"name":  2,
		"is":    1,
		"Payne": 1,
	}

	if !reflect.DeepEqual(res, m) {
		t.Errorf("TestLatin_simpleOneWord: Not equal")
	}
}
func TestLatin_simpleTwoWords(t *testing.T) {
	text := "Hello, my name is Max Payne."
	entity := "Max Payne"

	res := Latin(text, entity)
	m := map[string]float64{
		"Hello": 4,
		"my":    3,
		"name":  2,
		"is":    1,
	}

	if !reflect.DeepEqual(res, m) {
		t.Errorf("TestLatin_simpleTwoWords: Not equal")
	}
}

func TestLatin_complexTwoWords(t *testing.T) {
	text := "Max Payne. Hello, my name is Max Payne. I'm the best human and I'm Max. The real Max Payne, oh yes!"
	entity := "Max Payne"

	res := Latin(text, entity)
	m := map[string]float64{
		"Hello": 6.67,
		"my":    6.33,
		"name":  6,
		"is":    5.67,
		"I'm":   6.5,
		"the":   6,
		"best":  6.33,
		"human": 6.67,
		"and":   7,
		"Max":   7.67,
		"The":   8,
		"real":  8.33,
		"oh":    10.33,
		"yes":   11.33,
	}

	if !reflect.DeepEqual(res, m) {
		t.Errorf("TestLatin_complexTwoWords: Not equal")
	}
}
