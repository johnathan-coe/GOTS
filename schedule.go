package main

type schedule struct {
	node            *node
	processor       int
	startTime       int
	prev            *schedule
	nodes           int
	schedFinishTime int
}

type scheduleStack struct {
	top   *schedule
	under *scheduleStack
}

type Walk struct {
	next      *schedule
	scheduled []*schedule
	last      []*schedule
	procEnd   []int
}

func NewWalk(s *schedule, nodes, processors int) Walk {
	return Walk{
		next:      s,
		scheduled: make([]*schedule, nodes),
		last:      make([]*schedule, processors),
		procEnd:   make([]int, processors),
	}
}

// Walk one step
func (walk *Walk) walk() {
	sched := walk.next

	if sched == nil {
		return
	}

	walk.scheduled[sched.node.index] = sched
	if (walk.last[sched.processor]) == nil {
		walk.last[sched.processor] = sched
		walk.procEnd[sched.processor] = sched.startTime + sched.node.weight
	}

	walk.next = sched.prev
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

func (stack *scheduleStack) pop() (*scheduleStack, *schedule) {
	return stack.under, stack.top
}

func (stack *scheduleStack) push(item *schedule) *scheduleStack {
	return &scheduleStack{
		top:   item,
		under: stack,
	}
}
