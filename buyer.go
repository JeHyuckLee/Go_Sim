package main

import (
	"evsim_golang/executor"
	"evsim_golang/system"
	"fmt"
)

// Seeding
type buyer struct {
	executor *executor.BehaviorModelExecutor
	buy      int
	msg      *system.SysMessage
}

func AM_buyer(instance_time, destruct_time float64, name, engine_name string, buy int) *buyer {
	m := &buyer{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m
	m.buy = buy
	m.executor.Behaviormodel.Insert_state("IDLE", 15)
	m.executor.Behaviormodel.Insert_state("BUY", 1) //나중에 멤버에게 입력받아서 집어넣어야함
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("start")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("buy")

	return m
}

func (m *buyer) Ext_trans(port string, msg *system.SysMessage) {
	//파종이 필요하다고 요청이 옴
	// if port == "seeding" {
	// 	m.executor.Cancel_rescheduling()
	// 	m.executor.Cur_state = "SEEDING"
	// }
}

func (m *buyer) Output() *system.SysMessage {
	//가능한 수확량선에서 필요한 만큼 파종을 함
	fmt.Println("Buyer: buy a tomato : ", m.buy)
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "buy")
	msg.Insert(m.buy)
	return msg
}

func (m *buyer) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "BUY" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IDLE"
	}
}
