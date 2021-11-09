package main

import (
	"io/ioutil"
	"strconv"

	"github.com/goccy/go-graphviz"
)

type edge struct {
	weight int
	other  *node
}

type node struct {
	index   int
	name    string
	weight  int
	inEdges []edge
}

// Parse a graph from a .dot file
func parseGraph(path string) []*node {
	b, _ := ioutil.ReadFile(path)
	g, _ := graphviz.ParseBytes(b)

	// List of nodes in the graph
	nodes := make([]*node, g.NumberNodes())

	// Map of node names to nodes
	nodeMap := make(map[string]*node)

	// For each node
	n := 0
	for i := g.FirstNode(); i != nil; i = g.NextNode(i) {
		inEdges := make([]edge, g.Degree(i, 1, 0))

		e := 0
		for ed := g.FirstIn(i); ed != nil; ed = g.NextIn(ed) {
			weight, _ := strconv.Atoi(ed.Get("Weight"))

			inEdges[e] = edge{
				weight: weight,
				other:  nodeMap[ed.Node().Name()],
			}
			e++
		}

		weight, _ := strconv.Atoi(i.Get("Weight"))
		node := &node{
			index:   n,
			name:    i.Name(),
			weight:  weight,
			inEdges: inEdges,
		}

		nodes[n] = node
		nodeMap[i.Name()] = node
		n++
	}

	return nodes
}
