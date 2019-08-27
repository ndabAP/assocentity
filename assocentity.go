package assocentity

import (
	"math"

	"github.com/ndabAP/assocentity/v5/internal/iterator"
	"github.com/ndabAP/assocentity/v5/tokenize"
)

type direction int

var (
	posDir direction = 1
	negDir direction = -1
)

// Do returns the entity distances
func Do(tokenizer tokenize.Tokenizer, dps tokenize.PoSDetermer, entities []string) (map[tokenize.Token]float64, error) {
	determTok, err := dps.Determ(tokenizer)
	if err != nil {
		return map[tokenize.Token]float64{}, err
	}

	var distAccum = make(map[tokenize.Token][]float64)

	// Prepare for generic iterator
	di := make(iterator.Elements, len(determTok))
	for i, v := range determTok {
		di[i] = v
	}

	entityTokens, err := tokenizer.TokenizeEntities()
	if err != nil {
		return map[tokenize.Token]float64{}, err
	}

	determTokTraverser := iterator.New(di)
	for determTokTraverser.Next() {
		determTokIdx := determTokTraverser.CurrPos()

		// Check for entity
		var (
			dist            float64
			isEntity        bool
			entityTraverser *iterator.Iterator
		)
		// Check if entity
		for entityIdx := range entityTokens {
			// Prepare for generic iterator
			e := make(iterator.Elements, len(entityTokens[entityIdx]))
			for i, v := range entityTokens[entityIdx] {
				e[i] = v
			}

			entityTraverser = iterator.New(e)
			for entityTraverser.Next() {
				isEntity = determTokTraverser.CurrElem().(tokenize.Token).Token == entityTraverser.CurrElem().(tokenize.Token).Token
				// Check if first token matches the entity token
				if !isEntity {
					break
				}

				// Check for next token
				determTokTraverser.Next()
			}

			// Skip entity
			if isEntity {
				determTokTraverser.SetPos(determTokIdx + entityTraverser.Len() - 1)

				break
			}
		}

		if isEntity {
			continue
		}

		determTokTraverser.SetPos(determTokIdx)

		// Iterate positive direction
		posTraverser := iterator.New(di)
		posTraverser.SetPos(determTokIdx)
		for posTraverser.Next() {
			posTraverserIdx := posTraverser.CurrPos()
			if ok, len := isPartOfEntity(posTraverser, entityTokens, posDir); ok {
				distAccum[determTokTraverser.CurrElem().(tokenize.Token)] = append(distAccum[determTokTraverser.CurrElem().(tokenize.Token)], dist)

				// Skip about entity
				posTraverser.SetPos(posTraverserIdx + len - 1)

				continue
			}

			// Reset because isPartOfEntity is mutating
			posTraverser.SetPos(posTraverserIdx)

			dist++
		}

		// Iterate negative direction
		dist = 0
		negTraverser := iterator.New(di)
		negTraverser.SetPos(determTokIdx)
		for negTraverser.Prev() {
			negTraverserIdx := negTraverser.CurrPos()
			if ok, len := isPartOfEntity(negTraverser, entityTokens, negDir); ok {
				distAccum[determTokTraverser.CurrElem().(tokenize.Token)] = append(distAccum[determTokTraverser.CurrElem().(tokenize.Token)], dist)

				// Skip about entity
				negTraverser.SetPos(negTraverserIdx - len + 1)

				continue
			}

			// Reset because isPartOfEntity is mutating
			negTraverser.SetPos(negTraverserIdx)

			dist++
		}
	}

	assoccEntities := make(map[tokenize.Token]float64)
	// Calculate the distances
	for elem, dist := range distAccum {
		assoccEntities[elem] = avg(dist)
	}

	return assoccEntities, nil
}

// Iterates through entity tokens and returns true if found and positions to skip
func isPartOfEntity(determTokTraverser *iterator.Iterator, entityTokens [][]tokenize.Token, dir direction) (bool, int) {
	var (
		isEntity        bool
		entityTraverser *iterator.Iterator
	)
	for entityIdx := range entityTokens {
		// Prepare for generic iterator
		e := make(iterator.Elements, len(entityTokens[entityIdx]))
		for i, v := range entityTokens[entityIdx] {
			e[i] = v
		}

		entityTraverser = iterator.New(e)

		if dir == posDir {
			// Positive direction
			for entityTraverser.Next() {
				isEntity = determTokTraverser.CurrElem().(tokenize.Token).Token == entityTraverser.CurrElem().(tokenize.Token).Token
				// Check if first token matches the entity token
				if !isEntity {
					break
				}

				// Check for next token
				determTokTraverser.Next()
			}

		} else if dir == negDir {
			// Negative direction
			entityTraverser.SetPos(entityTraverser.Len() - 1)
			for entityTraverser.Prev() {
				isEntity = determTokTraverser.CurrElem().(tokenize.Token).Token == entityTraverser.CurrElem().(tokenize.Token).Token
				// Check if first token matches the entity token
				if !isEntity {
					break
				}

				// Check for next token
				determTokTraverser.Prev()
			}
		}

		if isEntity {
			return isEntity, entityTraverser.Len()
		}
	}

	return isEntity, entityTraverser.Len()
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
