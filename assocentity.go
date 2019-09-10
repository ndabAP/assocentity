package assocentity

import (
	"math"

	"github.com/ndabAP/assocentity/v6/internal/iterator"
	"github.com/ndabAP/assocentity/v6/tokenize"
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

	entityTokens, err := tokenizer.TokenizeEntities()
	if err != nil {
		return map[tokenize.Token]float64{}, err
	}

	var distAccum = make(map[tokenize.Token][]float64)

	// Prepare for generic iterator
	di := make(iterator.Elements, len(determTok))
	for i, v := range determTok {
		di[i] = v
	}

	determTokTraverser := iterator.New(di)
	for determTokTraverser.Next() {
		determTokIdx := determTokTraverser.CurrPos()

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
			for entityTraverser.Next() {
				isEntity = determTokTraverser.CurrElem().(tokenize.Token) == entityTraverser.CurrElem().(tokenize.Token)
				// Check if first token matches the entity token
				if !isEntity {
					break
				}

				// Check for next token
				determTokTraverser.Next()
			}

			if isEntity {
				break
			}
		}

		if isEntity {
			// Skip about entity positions
			determTokTraverser.SetPos(determTokIdx + entityTraverser.Len() - 1)

			continue
		}

		var dist float64

		// Reset because mutated
		determTokTraverser.SetPos(determTokIdx)

		// Iterate positive direction
		posTraverser := iterator.New(di)
		posTraverser.SetPos(determTokIdx)
		for posTraverser.Next() {
			posTraverserIdx := posTraverser.CurrPos()

			ok, len := isPartOfEntity(posTraverser, entityTokens, posDir)
			if ok {
				distAccum[determTokTraverser.CurrElem().(tokenize.Token)] = append(distAccum[determTokTraverser.CurrElem().(tokenize.Token)], dist)

				// Skip about entity
				posTraverser.SetPos(posTraverserIdx + len - 1)
			}

			dist++

			if ok {
				continue
			}

			// Reset because isPartOfEntity is mutating
			posTraverser.SetPos(posTraverserIdx)
		}

		// Iterate negative direction
		dist = 0
		negTraverser := iterator.New(di)
		negTraverser.SetPos(determTokIdx)
		for negTraverser.Prev() {
			negTraverserIdx := negTraverser.CurrPos()

			ok, len := isPartOfEntity(negTraverser, entityTokens, negDir)
			if ok {
				distAccum[determTokTraverser.CurrElem().(tokenize.Token)] = append(distAccum[determTokTraverser.CurrElem().(tokenize.Token)], dist)

				// Skip about entity
				negTraverser.SetPos(negTraverserIdx - len + 1)
			}

			dist++

			if ok {
				continue
			}

			// Reset because isPartOfEntity is mutating
			negTraverser.SetPos(negTraverserIdx)
		}
	}

	assocEntities := make(map[tokenize.Token]float64)
	// Calculate the distances
	for elem, dist := range distAccum {
		assocEntities[elem] = avg(dist)
	}

	return assocEntities, nil
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
