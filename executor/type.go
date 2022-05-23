package executor

import "evsim_golang/system"

type min_schedule_heap []*BehaviorModelExecutor

func (m min_schedule_heap) Len() int {
	return len(m)
}
func (m min_schedule_heap) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
func (m min_schedule_heap) Less(i, j int) bool {
	return m[i].Get_req_time() < m[j].Get_req_time()
}

func (m *min_schedule_heap) Push(elem interface{}) {
	*m = append(*m, elem.(*BehaviorModelExecutor))
}

func (m *min_schedule_heap) Pop() interface{} {
	old := *m
	n := len(old)
	elem := old[n-1]
	*m = old[0 : n-1]
	return elem
}

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
	return false
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
