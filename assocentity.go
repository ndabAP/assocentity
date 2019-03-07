// Package assocentity returns the average distance from words to a given entity.
package assocentity

import (
	"math"
)

// Romance accepts all Romance based languages.
func Romance(text, entity string) map[string]float64 {
	var textSplit []string
	var entitySplit []string
	textSplit = tokenize(text)
	entitySplit = tokenize(entity)

	// Find all entity positions as slice indices
	var entityPositions []int
	for i := range textSplit {
		if found := isSliceSubset(textSplit, entitySplit, i); found {
			for range entitySplit {
				entityPositions = append(entityPositions, i)
				i++
			}
		}
	}

	// If entity is not present, return early
	if len(entityPositions) == 0 {
		return map[string]float64{}
	}

	// Find word positions excluding entity
	wordPositions := make(map[string][]int)
	for i, word := range textSplit {
		found := false
		// Exclude entity
		for _, pos := range entityPositions {
			if i == pos {
				found = true
			}
		}

		if !found {
			wordPositions[word] = append(wordPositions[word], i)
		}
	}

	batchedEntityPositions := batch(entityPositions, len(entitySplit))
	distances := make(map[string][]float64)
	// Find distances between entity and words
	for word, positions := range wordPositions {
		for _, pos := range positions {
			for _, batch := range batchedEntityPositions {
				first := batch[0]
				last := batch[len(batch)-1]

				d := 0.0
				// Determinate relative entity position from word
				if first < pos {
					d = float64(pos - last)
				} else {
					d = float64(pos - first)
				}

				distances[word] = append(distances[word], math.Abs(d))
			}
		}
	}

	// Calculate the average distances
	assocentity := make(map[string]float64)
	for word, distances := range distances {
		assocentity[word] = round(avg(distances))
	}

	return assocentity
}

// Checks if next elements of given slice equals antoher given slice
func isSliceSubset(data, subset []string, index int) bool {
	hits := 0
	for i, sub := range subset {
		// Check for slice overflow
		if index+i > len(data)-1 {
			goto End
		}

		if data[index+i] == sub {
			hits++
		}
	}

End:
	return hits == len(subset)
}

// Create batches of given size
func batch(data []int, size int) [][]int {
	var batches [][]int
	for size < len(data) {
		data, batches = data[size:], append(batches, data[0:size:size])
	}

	batches = append(batches, data)

	return batches
}
