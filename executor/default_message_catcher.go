package executor

import (
	"evsim_golang/definition"
	"evsim_golang/system"
	"fmt"
)

type DefaultMessageCatcher struct {
	executor *BehaviorModelExecutor
}

func NewDMC(instance_time, destruct_time float64, name, engine_name string) *DefaultMessageCatcher {
	dmc := DefaultMessageCatcher{}
	dmc.executor = NewExecutor(instance_time, destruct_time, name, engine_name)
	dmc.executor.AbstractModel = &dmc
	dmc.executor.Init_state("IDLE")
	dmc.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	dmc.executor.Behaviormodel.CoreModel.Insert_input_port("uncauhth")
	return &dmc
}

func (d *DefaultMessageCatcher) Ext_trans(port string, msg *system.SysMessage) {
	fmt.Println("dmc ext_trans")
}
func (d *DefaultMessageCatcher) Int_trans() {
	fmt.Println("dmc inttrans")
}

func (d *DefaultMessageCatcher) Output() *system.SysMessage {
	msg := system.SysMessage{}
	fmt.Println("dmcoutput", msg.Get_dst())
	return &msg
}

func (d *DefaultMessageCatcher) Time_advance() float64 {
	return definition.Infinite
}
