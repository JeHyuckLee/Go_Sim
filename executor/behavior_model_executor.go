package executor

import (
	"evsim_golang/definition"
	"evsim_golang/model"
	"evsim_golang/system"
	"fmt"
	"math"
)

type BehaviorModelExecutor struct {
	sysobject     *system.SysObject
	Behaviormodel *model.Behaviormodel

	_cancel_reshedule_f bool //리스케쥴링펑션의 실행 여부
	engine_name         string
	Cur_state           string
	Instance_t          float64
	Destruct_t          float64
	Next_event_t        float64
	requestedTime       float64
}

func (b *BehaviorModelExecutor) String() string {
	return fmt.Sprintf("[N]:{%s}, [S]:{%s}", b.Behaviormodel.CoreModel.Get_name(), b.Cur_state)
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
	return b.Instance_t
}

func (b *BehaviorModelExecutor) Get_destruct_time() float64 {
	return b.Instance_t
}

func (b *BehaviorModelExecutor) Init_state(state string) {
	b.Cur_state = state
}

func (b *BehaviorModelExecutor) Ext_trans(port string, msg interface{}) {

}

func (b *BehaviorModelExecutor) Int_trans(port, msg string) {

}

func (b *BehaviorModelExecutor) Output() *system.SysMessage {
	var something *system.SysMessage
	return something
}
func (b *BehaviorModelExecutor) Time_advance() float64 {
	for key := range b.Behaviormodel.States {
		if key == b.Cur_state {
			return b.Behaviormodel.States[b.Cur_state]
		}
	}
	return -1
}
func (b *BehaviorModelExecutor) Set_req_time(global_time float64, elapsed_time int) {
	elapsed_time = 0
	if b.Time_advance() == definition.Infinite {
		b.Next_event_t = definition.Infinite
		b.requestedTime = definition.Infinite
	} else {
		if b._cancel_reshedule_f {
			b.requestedTime = math.Min(b.Next_event_t, global_time+b.Time_advance())
		} else {
			b.requestedTime = global_time + b.Time_advance()
		}
	}
}
func (b *BehaviorModelExecutor) Get_req_time() float64 {
	if b._cancel_reshedule_f {
		b._cancel_reshedule_f = false
	}
	b.Next_event_t = b.requestedTime
	return b.requestedTime
}

func NewExecutor(instantiate_time, destruct_time float64, name, engine_name string) *BehaviorModelExecutor {
	if instantiate_time == 0 {
		instantiate_time = math.Inf(1)
	}
	if destruct_time == 0 {
		destruct_time = math.Inf(1)
	}

	b := BehaviorModelExecutor{}
	b.engine_name = engine_name
	b.Instance_t = instantiate_time
	b.Destruct_t = destruct_time
	b.sysobject = system.NewSysObject()
	b.Behaviormodel = model.NewBehaviorModel(name)
	b.requestedTime = math.Inf(1)
	return &b
}
