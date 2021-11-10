package main

import (
	"os"
	"strconv"
)

func prune() bool {

	return false
}

func findOptimalSchedule(g []*node, processors int) *schedule {
	numNodes := len(g)

	// Schedule the first node on processor 0
	seed := &schedule{
		node:            g[0],
		processor:       0,
		startTime:       0,
		nodes:           1,
		schedFinishTime: g[0].weight,
	}

	stack := &scheduleStack{
		top:   seed,
		under: nil,
	}

	var best *schedule = nil
	for stack != nil {
		var n *schedule
		stack, n = stack.pop()

		// Ensure it is still a candidate for an optimal schedule
		if best != nil && n.schedFinishTime+n.node.criticalPath >= best.schedFinishTime {
			continue
		}

		// Nodes and the point in the schedule they were introduced at
		scheduled := make([]*schedule, len(g))
		finishTime := make([]int, processors)
		maxProc := 0

		// Walk schedules to get nodes scheduled
		for sched := n; sched != nil; sched = sched.prev {
			scheduled[sched.node.index] = sched

			if (finishTime[sched.processor]) == 0 {
				finishTime[sched.processor] = sched.startTime + sched.node.weight
				maxProc = max(maxProc, sched.processor)
			}
		}

		// Range of processors that we can schedule a new task on
		validProcessors := min(maxProc+2, processors)

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

			satisfiedAt := finishTime

			for _, dep := range s.inEdges {
				depNode := scheduled[dep.other.index]
				afterComms := depNode.startTime + depNode.node.weight + dep.weight

				for i := 0; i < validProcessors; i++ {
					if i != depNode.processor {
						satisfiedAt[i] = max(satisfiedAt[i], afterComms)
					}
				}
			}

			for i := 0; i < validProcessors; i++ {
				start := satisfiedAt[i]
				schedFinish := max(n.schedFinishTime, start+s.weight)

				// This is an optimal candidate
				if best == nil || schedFinish+s.criticalPath < best.schedFinishTime {
					nodes := n.nodes + 1

					newSched := &schedule{
						node:            s,
						processor:       i,
						startTime:       start,
						prev:            n,
						nodes:           nodes,
						schedFinishTime: schedFinish,
					}

					// It is complete
					if nodes == numNodes {
						best = newSched
					} else {
						stack = stack.push(newSched)
					}
				}
			}
		}
	}

	return best
}

func main() {
	// Parse args
	path := os.Args[1]
	processors, _ := strconv.Atoi(os.Args[2])

	nodes := parseGraph(path)
	s := findOptimalSchedule(nodes, processors)
	println(s.schedFinishTime)

	// Walk schedules to get nodes scheduled
	sched := s
	for sched != nil {
		println(sched.node.name, sched.startTime, sched.processor)
		sched = sched.prev
	}
}
