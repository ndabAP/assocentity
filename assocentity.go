// Package assocentity returns the average distance from words to a given entity.
package assocentity

import (
	"errors"
	"math"

	"gopkg.in/jdkato/prose.v2"
)

type tokenized []string

var errorNoEntityFound = errors.New("no entity was found inside text")

// Make accepts a text, entities including aliases and a tokenizer which defaults to an English tokenizer.
func Make(text string, entities []string, tokenizer func(string) ([]string, error)) (map[string]float64, error) {
	var tokenizedText tokenized
	var tokenizedEntities []tokenized
	var err error
	// Apply user given tokenizer if not nil
	if tokenizer == nil {
		tokenizedText, err = englishTokenizer(text)

		var tokenizedEntity tokenized
		for _, entity := range entities {
			tokenizedEntity, err = englishTokenizer(entity)
			if err != nil {
				return nil, err
			}

			tokenizedEntities = append(tokenizedEntities, tokenizedEntity)
		}
		if err != nil {
			return nil, err
		}
	} else {
		tokenizedText, err = tokenizer(text)
		if err != nil {
			return nil, err
		}

		var tokenizedEntity tokenized
		for _, entity := range entities {
			tokenizedEntity, err = tokenizer(entity)
			if err != nil {
				return nil, err
			}

			tokenizedEntities = append(tokenizedEntities, tokenizedEntity)
		}
	}

	entityPositions := findEntityPositions(tokenizedText, tokenizedEntities)
	if len(entityPositions) == 0 {
		return nil, errorNoEntityFound
	}
	wordPositions := findWordPositions(tokenizedText, entityPositions)
	wordsDistances := findWordEntityDistances(wordPositions, entityPositions)

	weighting := make(map[string]float64)
	for word, distances := range wordsDistances {
		avg := average(distances)
		weighting[word] = avg
	}

	return weighting, nil
}

// Tokenizes English words.
func englishTokenizer(text string) (tokenized, error) {
	document, err := prose.NewDocument(string(text))
	if err != nil {
		return nil, err
	}

	var tokenizedText tokenized
	for _, token := range document.Tokens() {
		tokenizedText = append(tokenizedText, token.Text)
	}

	return tokenizedText, nil
}

// Returns entity indices including aliases.
func findEntityPositions(tokenizedText tokenized, tokenizedEntites []tokenized) [][]int {
	var entityPositions [][]int
	hits := 0
	for _, tokenizedEntity := range tokenizedEntites {
	TokenizedTextLoop:
		for i := range tokenizedText {
			if found := isSliceSubset(tokenizedText, tokenizedEntity, i); found {
				// Check if entity position already found
				for _, entityPosition := range entityPositions {
					for _, position := range entityPosition {
						if position == i {
							continue TokenizedTextLoop
						}
					}
				}

				entityPositions = append(entityPositions, []int{})
				for range tokenizedEntity {
					entityPositions[hits] = append(entityPositions[hits], i)
					i++
				}
				hits++
			}

		}
	}

	return entityPositions
}

// Returns distances from entity to words.
func findWordEntityDistances(wordPositions map[string][]int, entityPositions [][]int) map[string][]float64 {
	wordDistances := make(map[string][]float64)
	for word, positions := range wordPositions {
		for _, wordPosition := range positions {
			for _, entityPosition := range entityPositions {
				firstEntityPosition := entityPosition[0]
				lastEntityPosition := entityPosition[len(entityPosition)-1]

				distance := 0.0
				if firstEntityPosition < wordPosition {
					distance = math.Abs(float64(wordPosition - lastEntityPosition))
				} else {
					distance = math.Abs(float64(wordPosition - firstEntityPosition))
				}

				wordDistances[word] = append(wordDistances[word], distance)
			}
		}
	}

	return wordDistances
}

// Returns word indices.
func findWordPositions(tokenizedText tokenized, entityPositions [][]int) map[string][]int {
	wordPositions := make(map[string][]int)
	for i, word := range tokenizedText {
		found := false
		for _, entityPosition := range entityPositions {
			if isInSlice(i, entityPosition) {
				found = true
			}
		}

		if !found {
			wordPositions[word] = append(wordPositions[word], i)
		}
	}

	return wordPositions
}

// Checks if next elements of given slice equals antoher given slice.
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

// Returns "true" if integer is in slice, else "false".
func isInSlice(n int, sl []int) bool {
	for _, e := range sl {
		if e == n {
			return true
		}
	}

	return false
}

// Returns the average of a float slice.
func average(xs []float64) float64 {
	total := 0.0
	for _, v := range xs {
		total += v
	}

	return round(total / float64(len(xs)))
}

// Rounds to nearest 0.5.
func round(x float64) float64 {
	return math.Round(x/0.5) * 0.5
}
