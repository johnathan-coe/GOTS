package main

import (
	"os"
	"strconv"
)

func prune(s *schedule, precededByPrev bool) bool {
	return !precededByPrev && s.processor < s.prev.processor
}

func FindOptimalSchedule(g []*node, processors int) *schedule {
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

		walk := &partialWalk{
			next:       n,
			scheduled:  make([]*schedule, len(g)),
			lastOnProc: make([]*schedule, processors),
			procEnd:    make([]int, processors),
		}

		for index, s := range g {
			if walk.walkTillIndex(index) != nil {
				continue
			}

			depsSatisfied := true
			for _, dep := range s.inEdges {
				if walk.walkTillIndex(dep.other.index) == nil {
					depsSatisfied = false
					break
				}
			}

			if !depsSatisfied {
				continue
			}

			// Resolve other processors till we hit one with no nodes
			validProcessors := 1
			for ; validProcessors < processors; validProcessors++ {
				if walk.walkTillProc(validProcessors-1) == nil {
					break
				}
			}

			satisfiedAt := walk.procEnd

			// We've already walked deps
			for _, dep := range s.inEdges {
				depNode := walk.scheduled[dep.other.index]
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
					newSched := &schedule{
						node:            s,
						processor:       i,
						startTime:       start,
						prev:            n,
						nodes:           n.nodes + 1,
						schedFinishTime: schedFinish,
					}

					// It is complete
					if newSched.nodes == numNodes {
						best = newSched
					} else {
						if !prune(newSched, n.node.goesTo[s.index]) {
							stack = stack.push(newSched)
						}
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
	s := FindOptimalSchedule(nodes, processors)
	println(s.schedFinishTime)

	// Walk schedules to get nodes scheduled
	sched := s
	for sched != nil {
		println(sched.node.name, sched.startTime, sched.processor)
		sched = sched.prev
	}
}
