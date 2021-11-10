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

type partialWalk struct {
	next       *schedule
	scheduled  []*schedule
	lastOnProc []*schedule
	procEnd    []int
}

func (walk *partialWalk) walk() {
	sched := walk.next

	if sched == nil {
		return
	}

	walk.scheduled[sched.node.index] = sched
	if (walk.lastOnProc[sched.processor]) == nil {
		walk.lastOnProc[sched.processor] = sched
		walk.procEnd[sched.processor] = sched.startTime + sched.node.weight
	}

	walk.next = sched.prev
}

func (walk *partialWalk) walkTillProc(p int) *schedule {
	for walk.lastOnProc[p] == nil && walk.next != nil {
		walk.walk()
	}

	return walk.lastOnProc[p]
}

func (walk *partialWalk) walkTillIndex(i int) *schedule {
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
