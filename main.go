package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/oleiade/lane"
)

type schedule struct {
	node       *cgraph.Node
	processor  int
	startTime  int
	prev       *schedule
	nodes      int
	finishTime int
}

func nodeID(n *cgraph.Node) int {
	id, err := strconv.Atoi(n.Name())

	if err != nil {
		log.Fatal("Node name invalid!")
	}

	return id
}

func edgeWeight(e *cgraph.Edge) int {
	weight, err := strconv.Atoi(e.Get("Weight"))

	if err != nil {
		log.Fatal("Cannot get weight for edge")
	}

	return weight
}

func nodeWeight(e *cgraph.Node) int {
	weight, err := strconv.Atoi(e.Get("Weight"))

	if err != nil {
		log.Fatal("Cannot get weight for node")
	}

	return weight
}

func Max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func findOptimalSchedule(g *cgraph.Graph, processors int) *schedule {
	scheduleStack := lane.NewStack()

	// Schedule the first node on processor 1
	seed := schedule{
		node:       g.FirstNode(),
		processor:  0,
		startTime:  0,
		nodes:      1,
		finishTime: nodeWeight(g.FirstNode()),
	}
	scheduleStack.Push(&seed)

	var best *schedule = nil

	for !scheduleStack.Empty() {
		n := scheduleStack.Pop().(*schedule)

		if n.nodes == g.NumberNodes() {
			if best == nil || n.finishTime < best.finishTime {
				best = n
			}
		}

		// Nodes and the point in the schedule they were introduced at
		scheduled := make([]*schedule, g.NumberNodes())
		earliestStart := make([]int, processors)

		// Walk schedules to get nodes scheduled
		sched := n
		for sched != nil {
			scheduled[nodeID(sched.node)] = sched
			earliestStart[sched.processor] = Max(sched.finishTime, earliestStart[sched.processor])
			sched = sched.prev
		}

		// For all unscheduled
		for s := g.FirstNode(); s != nil; s = g.NextNode(s) {
			if scheduled[nodeID(s)] != nil {
				continue
			}

			depsSatisfied := true
			for dep := g.FirstIn(s); dep != nil; dep = g.NextIn(dep) {
				if scheduled[nodeID(dep.Node())] == nil {
					depsSatisfied = false
					break
				}
			}

			if !depsSatisfied {
				continue
			}

			for i := 0; i < processors; i++ {
				// Get the time all our dependencies are satisfied at
				satisfiedAt := 0
				for dep := g.FirstIn(s); dep != nil; dep = g.NextIn(dep) {
					pre := scheduled[nodeID(dep.Node())]

					if pre.processor != i {
						end := pre.startTime + nodeWeight(pre.node) + edgeWeight(dep)
						satisfiedAt = Max(satisfiedAt, end)
					}
				}

				start := Max(earliestStart[i], satisfiedAt)

				// Add new schedule to the stack
				scheduleStack.Push(&schedule{
					node:       s,
					processor:  i,
					startTime:  start,
					prev:       n,
					nodes:      n.nodes + 1,
					finishTime: Max(n.finishTime, start+nodeWeight(s)),
				})
			}
		}
	}

	return best
}

func main() {
	// Parse args
	path := os.Args[1]
	processors, err := strconv.Atoi(os.Args[2])

	if err != nil {
		log.Fatal("Failed to parse processor number")
	}

	// Read graph from file
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// Parse graph
	graph, err := graphviz.ParseBytes(b)

	s := findOptimalSchedule(graph, processors)
	println(s.finishTime)

	// Walk schedules to get nodes scheduled
	sched := s
	for sched != nil {
		println(sched.node.Name(), sched.startTime, sched.processor)
		sched = sched.prev
	}
}
