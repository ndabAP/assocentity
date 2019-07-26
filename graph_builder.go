package assocentity

import (
	"fmt"
	"math"

	"github.com/ndabAP/assocentity/v3/graph"
)

func buildGraph(tokens, entities []string) (map[string]float64, error) {
	g := graph.NewGraph(tokens)

	assoccentities := make(map[string][]float64)
	// Retreive the iteratee
	next := g.Iteratee()
	// Iterate over graph
	for next() {
		node := g.GetCurrNode()
		// Ignore entities
		if isInSlice(node.Node, entities) {
			continue
		}

		var dist float64
		// Iterate right way
		for g.Next() {
			dist++

			if isInSlice(g.GetCurrNode().Node, entities) {
				assoccentities[node.Node] = append(assoccentities[node.Node], dist)
			}
		}

		g.SetCurrNode(node)

		dist = 0
		// Iterate left way
		for g.Prev() {
			dist++

			if isInSlice(g.GetCurrNode().Node, entities) {
				assoccentities[node.Node] = append(assoccentities[node.Node], dist)
			}
		}

		g.SetCurrNode(node)
		next()

		fmt.Println(node.Node)
		// g.Next()
	}

	fmt.Println(assoccentities)

	// Calculate average word distances
	weighting := make(map[string]float64)
	for w, dist := range assoccentities {
		weighting[w] = avg(dist)
	}

	return weighting, nil
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

// Returns the avg of a float slice.
func avg(xs []float64) float64 {
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
