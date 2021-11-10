package main

func prune(s *schedule, precededByPrev bool) bool {
	if !precededByPrev {
		// This is most of our pruning
		if s.processor < s.prev.processor {
			return true
		}

		if s.processor == s.prev.processor && s.node.order < s.prev.node.order {
			return true
		}
	}

	return false
}

func FindOptimalSchedule(g *graph, processors int) *schedule {
	var stack *scheduleStack

	// Add all nodes with no dependencies
	for _, node := range g.nodes {
		if len(node.inEdges) == 0 {
			stack = stack.push(
				&schedule{
					node:            node,
					processor:       0,
					startTime:       0,
					nodes:           1,
					schedFinishTime: node.weight,
				})
		}
	}

	var best *schedule = nil
	for stack != nil {
		var n *schedule
		stack, n = stack.pop()

		// Ensure it is still a candidate for an optimal schedule
		if best != nil && n.schedFinishTime+n.node.criticalPath >= best.schedFinishTime {
			continue
		}

		// Make a new walk for this schedule "lazy loading"
		walk := NewWalk(n, g.numNodes, processors)

		for index, s := range g.nodes {
			if walk.scheduleForIndex(index) != nil {
				continue
			}

			depsSatisfied := true
			for _, dep := range s.inEdges {
				if walk.scheduleForIndex(dep.other.index) == nil {
					depsSatisfied = false
					break
				}
			}

			if !depsSatisfied {
				continue
			}

			// Resolve other processors till we hit one with no nodes
			validProcessors := processors
			for i := 0; i < processors; i++ {
				if walk.lastOnProc(i) == nil {
					validProcessors = i + 1
					break
				}
			}

			// The processors we care about are resolved
			satisfiedAt := make([]int, validProcessors)
			copy(satisfiedAt, walk.procEnd)

			// We've already walked deps
			for _, dep := range s.inEdges {
				depNode := walk.scheduled[dep.other.index]
				afterComms := depNode.startTime + depNode.node.weight + dep.weight

				for i := 0; i < validProcessors; i++ {
					if i != depNode.processor {
						if afterComms > satisfiedAt[i] {
							satisfiedAt[i] = afterComms
						}
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
					if newSched.nodes == g.numNodes {
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
