package assocentity

import (
	"reflect"
	"testing"
)

func TestRomance_Empty(t *testing.T) {
	text := "Hello world"
	entity := "Bye"

	res := Romance(text, entity)
	if !reflect.DeepEqual(res, map[string]float64{}) {
		t.Errorf("TestRomance_Empty: Not equal")
	}
}

func TestRomance_Subset(t *testing.T) {
	text := "Hello world"
	entity := "Helloworld"

	res := Romance(text, entity)
	if !reflect.DeepEqual(res, map[string]float64{}) {
		t.Errorf("TestRomance_Empty: Not equal")
	}
}

func TestRomance_Start(t *testing.T) {
	text := "Shang Tsung is my name."
	entity := "Shang Tsung"

	res := Romance(text, entity)
	m := map[string]float64{
		"is":   1,
		"my":   2,
		"name": 3,
	}
	if !reflect.DeepEqual(res, m) {
		t.Errorf("TestRomance_Empty: Not equal")
	}
}

func TestRomance_simpleOneWord(t *testing.T) {
	text := "Hello, my name is Max Payne."
	entity := "Max"

	res := Romance(text, entity)
	m := map[string]float64{
		"Hello": 4,
		"my":    3,
		"name":  2,
		"is":    1,
		"Payne": 1,
	}

	if !reflect.DeepEqual(res, m) {
		t.Errorf("TestRomance_simpleOneWord: Not equal")
	}
}
func TestRomance_simpleTwoWords(t *testing.T) {
	text := "Hello, my name is Max Payne."
	entity := "Max Payne"

	res := Romance(text, entity)
	m := map[string]float64{
		"Hello": 4,
		"my":    3,
		"name":  2,
		"is":    1,
	}

	if !reflect.DeepEqual(res, m) {
		t.Errorf("TestRomance_simpleTwoWords: Not equal")
	}
}

func TestRomance_simpleFourWords(t *testing.T) {
	text := `If you smell, what Dwayne "The Rock" Johnson is cooking?`
	entity := `Dwayne "The Rock" Johnson`

	res := Romance(text, entity)
	m := map[string]float64{
		"If":      4,
		"you":     3,
		"smell":   2,
		"what":    1,
		"is":      1,
		"cooking": 2,
	}

	if !reflect.DeepEqual(res, m) {
		t.Errorf("TestRomance_simpleTwoWords: Not equal")
	}
}

func TestRomance_complexTwoWords(t *testing.T) {
	text := "Shao Kahn is the embodiment of evil. Shao Kahn is easily recognizable by his intimidating stature."
	entity := "Shao Kahn"

	res := Romance(text, entity)
	m := map[string]float64{
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
	}

	if !reflect.DeepEqual(res, m) {
		t.Errorf("TestRomance_complexTwoWords: Not equal")
	}
}
