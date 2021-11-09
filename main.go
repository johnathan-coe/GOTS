package main

import (
	"log"
	"os"
	"strconv"
)

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func findOptimalSchedule(g []*node, processors int) *schedule {
	// Schedule the first node on processor 0
	seed := &schedule{
		node:       g[0],
		processor:  0,
		startTime:  0,
		nodes:      1,
		finishTime: g[0].weight,
	}

	var best *schedule = nil
	stack := &scheduleStack{
		top:   seed,
		under: nil,
	}

	for stack != nil {
		var n *schedule
		stack, n = stack.pop()

		// When we find a complete schedule
		if n.nodes == len(g) {
			if best == nil || n.finishTime < best.finishTime {
				best = n
			}
		}

		// Nodes and the point in the schedule they were introduced at
		scheduled := make([]*schedule, len(g))
		earliestStart := make([]int, processors)

		// Walk schedules to get nodes scheduled
		for sched := n; sched != nil; sched = sched.prev {
			scheduled[sched.node.index] = sched

			if (earliestStart[sched.processor]) == 0 {
				earliestStart[sched.processor] = sched.finishTime
			}
		}

		for index, s := range g {
			if scheduled[index] != nil {
				continue
			}

			depsSatisfied := true
			for _, dep := range s.inEdges {
				if scheduled[dep.other.index] == nil {
					depsSatisfied = false
					break
				}
			}

			if !depsSatisfied {
				continue
			}

			// If we can schedule, try different processors
			encounteredEmpty := false
			for i := 0; i < processors; i++ {
				// Two empty processors means we can skip
				empty := earliestStart[i] == 0
				if encounteredEmpty && empty {
					break
				}

				encounteredEmpty = empty || encounteredEmpty

				// Get the time all our dependencies are satisfied at
				satisfiedAt := 0
				for _, dep := range s.inEdges {
					pre := scheduled[dep.other.index]

					if pre.processor != i {
						end := pre.startTime + pre.node.weight + dep.weight
						satisfiedAt = max(satisfiedAt, end)
					}
				}

				start := max(earliestStart[i], satisfiedAt)
				finish := max(n.finishTime, start+s.weight)

				if best == nil || finish+s.criticalPath < best.finishTime {
					stack = stack.push(&schedule{
						node:       s,
						processor:  i,
						startTime:  start,
						prev:       n,
						nodes:      n.nodes + 1,
						finishTime: finish,
					})
				}
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

	nodes := parseGraph(path)
	s := findOptimalSchedule(nodes, processors)

	println(s.finishTime)

	// Walk schedules to get nodes scheduled
	sched := s
	for sched != nil {
		println(sched.node.name, sched.startTime, sched.processor)
		sched = sched.prev
	}
}
