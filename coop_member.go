package main

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
	"evsim_golang/system"
	"fmt"
)

// CoopMember
type cm_coopMember struct {
	am_seed    *coopMember_seed
	am_harvest *coopMember_harvest
	am_ship    *coopMember_ship
}

func CM_coopMember(instance_time, destruct_time float64, name, engine_name string, area int, harvest int, period int) *cm_coopMember {

	cell := cm_coopMember{}
	cell.am_seed = AM_seed(instance_time, destruct_time, name, engine_name, area, harvest)
	cell.am_harvest = AM_harvest(instance_time, destruct_time, name, engine_name, area, harvest, period)
	cell.am_ship = AM_ship(instance_time, destruct_time, name, engine_name)

	return &cell
}

// Seeding
type coopMember_seed struct {
	executor *executor.BehaviorModelExecutor
	area     int
	harvest  int
	msg      *system.SysMessage
}

func AM_seed(instance_time, destruct_time float64, name, engine_name string, area int, harvest int) *coopMember_seed {
	m := &coopMember_seed{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	//infor
	m.area = area
	m.harvest = harvest

	//statef
	m.executor.Behaviormodel.Insert_state("IDLE", 50)
	m.executor.Behaviormodel.Insert_state("SEEDING", 1) //나중에 멤버에게 입력받아서 집어넣어야함
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("seeding")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("harvest")

	return m
}

func (m *coopMember_seed) Ext_trans(port string, msg *system.SysMessage) {
	//파종이 필요하다고 요청이 옴
	if port == "seeding" {
		m.executor.Cancel_rescheduling()
		fmt.Println("State: Seeding in")
		m.executor.Cur_state = "SEEDING"
	}
}

func (m *coopMember_seed) Output() *system.SysMessage {
	//가능한 수확량선에서 필요한 만큼 파종을 함

	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "harvest")

	return msg
}

func (m *coopMember_seed) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "SEEDING" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IDLE"
	}
}

// Harvest
type coopMember_harvest struct {
	executor *executor.BehaviorModelExecutor
	area     int
	harvest  int
	period   int
	tomato   tomato
	msg      *system.SysMessage
}

func AM_harvest(instance_time, destruct_time float64, name, engine_name string, area int, harvest int, period int) *coopMember_harvest {
	m := &coopMember_harvest{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	//infor
	m.period = period
	m.area = area
	m.harvest = harvest
	//state
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("HARVEST", 1)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("harvest")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("shipment")

	return m
}

func (m *coopMember_harvest) Ext_trans(port string, msg *system.SysMessage) {
	//player가 해당 셀에 왔음
	if port == "harvest" {
		fmt.Println("[Seeding] => [Harvest]")
		m.executor.Cancel_rescheduling()

		m.tomato = tomato{m.harvest, m.period}
		m.executor.Cur_state = "HARVEST"
	}

}

func (m *coopMember_harvest) Output() *system.SysMessage {
	fmt.Println("Harvest...")
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "shipment")
	msg.Insert(m.tomato)

	return msg
}

func (m *coopMember_harvest) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "HARVEST" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IDLE"
	}
}

// Shipment
type coopMember_ship struct {
	executor *executor.BehaviorModelExecutor
	shipment int
	tomato   tomato
	msg      *system.SysMessage
}

func AM_ship(instance_time, destruct_time float64, name, engine_name string) *coopMember_ship {
	m := &coopMember_ship{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	//infor

	//statef
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("SHIPMENT", 1)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("shipment")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("in")

	return m
}

func (m *coopMember_ship) Ext_trans(port string, msg *system.SysMessage) {
	//player가 해당 셀에 왔음
	if port == "shipment" {
		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		m.tomato = data[0].(tomato)

		m.executor.Cur_state = "SHIPMENT"
	}
}

func (m *coopMember_ship) Output() *system.SysMessage {
	//check 에게 출력을 보내서 동작시킴
	fmt.Println("member Shipment quantity: ", m.tomato.Quantity)
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "in")
	msg.Insert(m.tomato)

	return msg
}

func (m *coopMember_ship) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "SHIPMENT" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IDLE"
	}
}
