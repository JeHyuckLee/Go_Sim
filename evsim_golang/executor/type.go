package executor

import (
	"evsim_golang/system"
)

const epsilon = 1e-14

type Object struct {
	object *BehaviorModelExecutor
	port   string
}

type i_event_queue struct {
	time float64
	msg  *system.SysMessage
}

type o_event_queue struct {
	time     float64
	msg_list interface{}
}

type input_heap []i_event_queue

func (eq input_heap) Len() int {
	return len(eq)
}

func (eq input_heap) Less(i, j int) bool {
	return eq[i].time < eq[j].time
}

func (eq input_heap) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
}

func (eq *input_heap) Push(elem interface{}) {
	*eq = append(*eq, elem.(i_event_queue))
}

func (eq *input_heap) Pop() interface{} {
	old := *eq
	n := len(old)
	elem := old[n-1]
	*eq = old[0 : n-1]

	return elem
}

//고루틴
