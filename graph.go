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
	index    int
	name     string
	weight   int
	inEdges  []edge
	outEdges []edge
	// Heuristics
	criticalPath int
	order        int
	goesTo       []bool
}

func (n *node) calculateCriticalPath() int {
	if n.criticalPath == 0 {
		for _, e := range n.outEdges {
			n.criticalPath = max(e.other.calculateCriticalPath()+e.other.weight, n.criticalPath)
		}
	}

	return n.criticalPath
}

func markOrder(g []*node) {
	readyQueue := make([]*node, 0)
	in := make([]int, len(g))

	for _, n := range g {
		incoming := len(n.inEdges)

		if incoming == 0 {
			readyQueue = append(readyQueue, n)
		}

		in[n.index] = incoming
	}

	order := 0
	for len(readyQueue) > 0 {
		ready := readyQueue[0]
		readyQueue = readyQueue[1:]

		ready.order = order
		order++

		for _, e := range ready.outEdges {
			child := e.other

			in[child.index]--

			if in[child.index] == 0 {
				readyQueue = append(readyQueue, child)
			}
		}
	}
}

// Parse a graph from a .dot file
func parseGraph(path string) []*node {
	b, _ := ioutil.ReadFile(path)
	g, _ := graphviz.ParseBytes(b)

	// List of nodes in the graph
	nodes := make([]*node, g.NumberNodes())

	// Map of node names to nodes
	nodeMap := make(map[string]*node)

	n := 0
	for i := g.FirstNode(); i != nil; i = g.NextNode(i) {
		weight, _ := strconv.Atoi(i.Get("Weight"))
		node := &node{
			index:  n,
			name:   i.Name(),
			weight: weight,
		}

		nodes[n] = node
		nodeMap[i.Name()] = node
		n++
	}

	for i := g.FirstNode(); i != nil; i = g.NextNode(i) {
		node := nodeMap[i.Name()]

		// Ingoing edges
		{
			node.inEdges = make([]edge, g.Degree(i, 1, 0))
			e := 0
			for ed := g.FirstIn(i); ed != nil; ed = g.NextIn(ed) {
				weight, _ := strconv.Atoi(ed.Get("Weight"))
				node.inEdges[e] = edge{
					weight: weight,
					other:  nodeMap[ed.Node().Name()],
				}
				e++
			}
		}

		// Outgoing edges
		{
			node.outEdges = make([]edge, g.Degree(i, 0, 1))
			e := 0
			for ed := g.FirstOut(i); ed != nil; ed = g.NextOut(ed) {
				weight, _ := strconv.Atoi(ed.Get("Weight"))
				node.outEdges[e] = edge{
					weight: weight,
					other:  nodeMap[ed.Node().Name()],
				}
				e++
			}
		}
	}

	for _, n := range nodes {
		n.calculateCriticalPath()

		// Bitfield of outgoing edges from this node
		n.goesTo = make([]bool, len(nodes))
		for _, out := range n.outEdges {
			n.goesTo[out.other.index] = true
		}
	}

	markOrder(nodes)

	return nodes
}
