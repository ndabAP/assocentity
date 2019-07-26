package graph

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/multi"
)

// Node represents a Node
type Node struct {
	Node string // Node
	id   int64  // Index to get unique hash
}

// ID returns the unique identifier for a graph node
func (n Node) ID() int64 {
	return n.id
}

// Graph represents a graph
type Graph struct {
	*multi.DirectedGraph          // Directed graph
	nodes                []string // Nodes
	iterator                      // To iterate over the graph
}

// Iterates over a graph
type iterator struct {
	currNode graph.Node // Current node
	id       int64      // Internal counter
}

// NewGraph returns a new graph
func NewGraph(nodes []string) Graph {
	// Create a new directed graph
	g := Graph{
		multi.NewDirectedGraph(),
		nodes,
		iterator{
			Node{},
			0,
		},
	}

	// Iterate over tokens
	for i, node := range nodes {
		// Add a node for every token
		n := Node{
			Node: node,
			id:   int64(i),
		}

		g.AddNode(n)

		// The first node is the current node
		if i == 0 {
			g.currNode = n
		}

		// Start adding lines at the second element
		if i > 0 {
			prev := nodes[i-1]
			line := g.NewLine(
				Node{
					Node: prev,
					id:   int64(i - 1),
				},
				n,
			)

			g.SetLine(line)
		}
	}

	return g
}

// Iteratee returns a function to iterate over the graph
func (g *Graph) Iteratee() func() bool {
	return func() bool {
		if g.id == 0 {
			g.id++

			return true
		}

		return g.Next()
	}
}

// Next calls the next nodes with given starting point
func (g *Graph) Next() bool {
	// We reached the end
	if g.id == int64(len(g.nodes)-1) {
		return false
	}

	from := g.From(g.id)
	if from.Next() {
		g.id++
		g.currNode = from.Node()

		return true
	}

	return false
}

// Prev calls the previous nodes with given starting point
func (g *Graph) Prev() bool {
	// We reached the start
	if g.id == 0 {
		return false
	}

	to := g.To(g.id)
	if to.Next() {
		g.id--
		g.currNode = to.Node()

		return true
	}

	return false
}

// SetCurrNode sets the current node
func (g *Graph) SetCurrNode(n Node) {
	g.id = n.id
	g.currNode = n
}

// GetCurrNode gets the current node
func (g *Graph) GetCurrNode() Node {
	return g.currNode.(Node)
}
