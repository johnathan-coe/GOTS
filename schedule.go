package main

type schedule struct {
	node       *node
	processor  int
	startTime  int
	prev       *schedule
	nodes      int
	finishTime int
}

type scheduleStack struct {
	top   *schedule
	under *scheduleStack
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
