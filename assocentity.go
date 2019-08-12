package assocentity

import (
	"math"

	"github.com/ndabAP/assocentity/v5/internal/iterator"
	"github.com/ndabAP/assocentity/v5/tokenize"
)

// Do returns the entity distances
func Do(tokenizer tokenize.Tokenizer, dps tokenize.PoSDetermer, j tokenize.Joiner, entities []string) (map[string]float64, error) {
	joined, err := j.Join(dps, tokenizer)
	if err != nil {
		return map[string]float64{}, err
	}

	var distAccum = make(map[string][]float64)

	// Prepare for generic iterator
	ji := make(iterator.Elements, len(joined))
	for i, v := range joined {
		ji[i] = v
	}

	joinedTraverser := iterator.New(ji)
	for joinedTraverser.Next() {
		joinedIdx := joinedTraverser.CurrPos()
		// Ignore entities
		if isInSlice(joinedTraverser.CurrElem().(string), entities) {
			continue
		}

		var dist float64

		// Iterate positive direction
		posTraverser := iterator.New(ji)
		posTraverser.SetPos(joinedIdx)
		for posTraverser.Next() {
			if isInSlice(posTraverser.CurrElem().(string), entities) {
				distAccum[joinedTraverser.CurrElem().(string)] = append(distAccum[joinedTraverser.CurrElem().(string)], dist)
			}

			dist++
		}

		// Iterate negative direction
		dist = 0
		negTraverser := iterator.New(ji)
		negTraverser.SetPos(joinedIdx)
		for negTraverser.Prev() {
			if isInSlice(negTraverser.CurrElem().(string), entities) {
				distAccum[joinedTraverser.CurrElem().(string)] = append(distAccum[joinedTraverser.CurrElem().(string)], dist)
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
