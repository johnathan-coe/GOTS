package main

import (
	"os"
	"strconv"
)

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
