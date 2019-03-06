package assocentity

import (
	"math"
	"strings"
)

// Latin returns the average distance from words to a given entity
func Latin(text, entity string) map[string]float64 {
	// Only allow English alphabet, ' and space
	var legalChars []rune
	for _, char := range text {
		if (char >= unicodecapa && char <= unicodecapz) ||
			(char >= unicodesma && char < unicodesmz) ||
			char == unicodespace ||
			char == unicodeapostrophe {
			legalChars = append(legalChars, char)
		}
	}

	// Split by empty space
	splittedText := strings.Split(string(legalChars), defaultsplit)
	splittedEntity := strings.Split(entity, defaultsplit)

	// Find entity positions as slice indices
	var entityPositions []int
	for i := range splittedText {
		found := scanLeftToRight(splittedText, splittedEntity, i)
		if found {
			for range splittedEntity {
				entityPositions = append(entityPositions, i)
				i++
			}
		}
	}

	// Find word positions excluding entity
	wordPositions := make(map[string][]int)
	for i, word := range splittedText {
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

	batched := batch(entityPositions, len(splittedEntity))
	// Find distances between entity and words
	distances := make(map[string][]float64)
	for word, positions := range wordPositions {
		for _, pos := range positions {
			for _, batch := range batched {
				first := batch[0]
				last := batch[len(batch)-1]

				d := 0.0
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
	res := make(map[string]float64)
	for word, distances := range distances {
		res[word] = round(avg(distances))
	}

	return res
}

// Checks if next elements of given slice equals antoher given slice
func scanLeftToRight(data, subset []string, index int) bool {
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

// Average from float slice
func avg(x []float64) float64 {
	total := 0.0
	for _, v := range x {
		total += v
	}

	return total / float64(len(x))
}

// Round two decimal places
func round(x float64) float64 {
	return math.Round(x*100) / 100
}
