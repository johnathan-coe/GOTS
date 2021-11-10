package main

type schedule struct {
	node            *node
	processor       int
	startTime       int
	prev            *schedule
	nodes           int
	schedFinishTime int
}

type sliceScheduleStack []*schedule

func (stack *sliceScheduleStack) push(s *schedule) {
	*stack = append(*stack, s)
}

func (stack *sliceScheduleStack) pop() *schedule {
	index := len(*stack) - 1
	element := (*stack)[index]
	*stack = (*stack)[:index]
	return element
}

type Walk struct {
	next      *schedule
	scheduled []*schedule
	last      []*schedule
}

func NewWalk(s *schedule, nodes, processors int) Walk {
	return Walk{
		next:      s,
		scheduled: make([]*schedule, nodes),
		last:      make([]*schedule, processors),
	}
}

// Walk one step
func (walk *Walk) walk() {
	sched := walk.next
	if sched == nil {
		return
	}

	walk.next = sched.prev

	walk.scheduled[sched.node.index] = sched
	if (walk.last[sched.processor]) == nil {
		walk.last[sched.processor] = sched
	}
}

func (walk *Walk) lastOnProc(p int) *schedule {
	for walk.last[p] == nil && walk.next != nil {
		walk.walk()
	}

	return walk.last[p]
}

func (walk *Walk) scheduleForIndex(i int) *schedule {
	for walk.scheduled[i] == nil && walk.next != nil {
		walk.walk()
	}

	return walk.scheduled[i]
}
