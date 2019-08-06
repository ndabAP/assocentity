package assocentity

import (
	"math"

	"github.com/ndabAP/assocentity/v3/internal/iterator"
	"github.com/ndabAP/assocentity/v3/tokenize"
)

// Do returns the entity distances
func Do(j tokenize.Joiner, tokenizer tokenize.Tokenizer, entities []string) (map[string]float64, error) {
	err := j.Join(tokenizer)
	if err != nil {
		return nil, err
	}

	var distAccum = make(map[string][]float64)
	tok := j.Tokens()
	joinedTraverser := iterator.New(&tok)
	for joinedTraverser.Next() {
		// Ignore entities
		if isInSlice(joinedTraverser.CurrElem(), entities) {
			continue
		}

		var dist float64

		// Iterate positive direction
		posTraverser := iterator.New(&tok)
		posTraverser.SetPos(joinedTraverser.CurrPos())
		for posTraverser.Next() {
			if isInSlice(posTraverser.CurrElem(), entities) {
				distAccum[joinedTraverser.CurrElem()] = append(distAccum[joinedTraverser.CurrElem()], dist)
			}

			dist++
		}

		dist = 0

		// Iterate negative direction
		negTraverser := iterator.New(&tok)
		negTraverser.SetPos(joinedTraverser.CurrPos())
		for negTraverser.Prev() {
			if isInSlice(negTraverser.CurrElem(), entities) {
				distAccum[joinedTraverser.CurrElem()] = append(distAccum[joinedTraverser.CurrElem()], dist)
			}

			dist++
		}
	}

	assoccEntities := make(map[string]float64)
	// Calculate the distances
	for elem, dist := range distAccum {
		assoccEntities[elem] = avg(dist)
	}

	return assoccEntities, nil
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
