package graph

import (
	"fmt"
	"math"

	"github.com/ndabAP/assocentity/v3/graph"
)

func buildGraph(tokens, entities []string) (map[string]float64, error) {
	g := graph.NewGraph(tokens)

	var dist float64
	assoccentities := make(map[string][]float64)
	// Retreive the iteratee
	next := g.Iteratee()
	// Iterate over graph
	for next() {
		fmt.Println(g.GetCurrNode())
		node := g.GetCurrNode()
		// Ignore entities
		if isInSlice(node.Node, entities) {
			continue
		}

		dist = 0
		right := g
		right.SetCurrNode(node)
		// Iterate right way
		for right.Next() {
			dist++

			if isInSlice(right.GetCurrNode().Node, entities) {
				assoccentities[node.Node] = append(assoccentities[node.Node], dist)
			}
		}

		dist = 0
		left := g
		left.SetCurrNode(node)
		// Iterate left way
		for left.Prev() {
			dist++

			if isInSlice(left.GetCurrNode().Node, entities) {
				assoccentities[node.Node] = append(assoccentities[node.Node], dist)
			}
		}
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
