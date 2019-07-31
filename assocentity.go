package assocentity

import (
	"math"

	"github.com/ndabAP/assocentity/v3/internal/generator"
	"github.com/ndabAP/assocentity/v3/tokenize"
)

// Assoc returns the entity distances
func Assoc(mp tokenize.Joiner, tokenizer tokenize.Tokenizer, entities []string) (map[string]float64, error) {
	multiplexed, err := mp.Join(tokenizer)
	if err != nil {
		return nil, err
	}

	var distAccum = make(map[string][]float64)
	multiplexedTraverer := generator.New(multiplexed)
	for multiplexedTraverer.Next() {
		// Ignore entities
		if isInSlice(multiplexedTraverer.CurrElem(), entities) {
			continue
		}

		var dist float64

		// Iterate positive direction
		posTraverser := multiplexedTraverer
		for posTraverser.Next() {
			dist++

			if isInSlice(posTraverser.CurrElem(), entities) {
				distAccum[posTraverser.CurrElem()] = append(distAccum[posTraverser.CurrElem()], dist)
			}
		}

		// Iterate negative direction
		negTraverser := multiplexedTraverer
		for negTraverser.Prev() {
			dist++

			if isInSlice(posTraverser.CurrElem(), entities) {
				distAccum[posTraverser.CurrElem()] = append(distAccum[posTraverser.CurrElem()], dist)
			}
		}
	}

	assoccentities := make(map[string]float64)
	// Calculate the distances
	for elem, dist := range distAccum {
		assoccentities[elem] = avg(dist)
	}

	return assoccentities, nil
}

// Checks if string is in slice
func isInSlice(x string, y []string) bool {
	for _, v := range y {
		if v == x {
			return true
		}
	}

	return false
}

// Returns the average of a float slice
func avg(xs []float64) float64 {
	total := 0.0
	for _, v := range xs {
		total += v
	}

	return round(total / float64(len(xs)))
}

// Rounds to nearest 0.5
func round(x float64) float64 {
	return math.Round(x/0.5) * 0.5
}
