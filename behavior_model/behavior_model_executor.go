package behaviormodel

import system_object "evsim_golang/system"

type BehaviorModelExecutor struct {
	sysobject     *system_object.SysObject
	behaviormodel *Behaviormodel

	_cancel_reshedule_f bool
	engine_name         string
	_instance_t         float64
	_destruct_t         float64
	_cur_state          string
	_next_event_t       int
	requestedTime       float64
}

func (b *BehaviorModelExecutor) Cancel_rescheduling() {
	b._cancel_reshedule_f = true
}

func (b *BehaviorModelExecutor) Get_engine_name() string {
	return b.engine_name
}

func (b *BehaviorModelExecutor) Set_engine_name(name string) {
	b.engine_name = name
}

func (b *BehaviorModelExecutor) Get_create_time() float64 {
	return b._instance_t
}

func (b *BehaviorModelExecutor) Get_destruct_time() float64 {
	return b._destruct_t
}

func (b *BehaviorModelExecutor) Init_state(state string) {
	b._cur_state = state
}

func (b *BehaviorModelExecutor) Ext_trans(port, msg string) {

}

func (b *BehaviorModelExecutor) Int_trans(port, msg string) {

}
func (b *BehaviorModelExecutor) Output() {

}

func (b *BehaviorModelExecutor) Time_advance(port, msg string) {

}
func (b *BehaviorModelExecutor) Set_req_time(port, msg string) {

}
func (b *BehaviorModelExecutor) Get_req_time(port, msg string) {

}

func NewExecutor(instantiate_time, destruct_time float64, name, engine_name string) *BehaviorModelExecutor {
	b := BehaviorModelExecutor{}
	b.engine_name = engine_name
	b._instance_t = instantiate_time
	b._destruct_t = destruct_time
	b.sysobject = system_object.NewSysObject()
	b.behaviormodel = NewBehaviorModel(name)
	return &b
}
