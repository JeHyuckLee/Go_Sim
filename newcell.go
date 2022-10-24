package main

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
	"evsim_golang/system"
	"fmt"
)

//CoopMember
type cm_coopMember struct {
	am_seed    *coopMember_seed
	am_harvest *coopMember_harvest
	am_ship    *coopMember_ship
}

func CM_coopMember(instance_time, destruct_time float64, name, engine_name string, area int) *cm_coopMember {

	cell := cm_coopMember{}
	cell.am_seed = AM_seed(instance_time, destruct_time, name, engine_name, area)
	cell.am_harvest = AM_harvest(instance_time, destruct_time, name, engine_name, area)
	cell.am_ship = AM_ship(instance_time, destruct_time, name, engine_name)

	return &cell
}

//Seeding
type coopMember_seed struct {
	executor *executor.BehaviorModelExecutor
	area     int
	harvest  int
	req_item interface{}
	msg      *system.SysMessage
}

func AM_seed(instance_time, destruct_time float64, name, engine_name string, area int) *coopMember_seed {
	m := &coopMember_seed{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	//infor
	m.area = area
	m.harvest = 3

	//statef
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("IN", 0)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("in")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("harvest")

	return m
}

func (m *coopMember_seed) Ext_trans(port string, msg *system.SysMessage) {
	//파종이 필요하다고 요청이 옴
	if port == "in" {
		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		m.req_item = data[0]

		m.executor.Cur_state = "IN"
	}
}

func (m *coopMember_seed) Output() *system.SysMessage {
	//가능한 수확량선에서 필요한 만큼 파종을 함
	fmt.Println("State: Seeding in")
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "harvest")

	//파종하기
	req_item := m.req_item.(int)
	possHarvest := m.area * m.harvest
	var seeding []int
	seeding = append(seeding, possHarvest)
	if possHarvest < req_item {
		restHarvest := req_item - possHarvest
		seeding = append(seeding, restHarvest)
	}

	msg.Insert(seeding)
	return msg
}

func (m *coopMember_seed) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "IN" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IDLE"
	}
}

//Harvest
type coopMember_harvest struct {
	executor *executor.BehaviorModelExecutor
	area     int
}

func AM_harvest(instance_time, destruct_time float64, name, engine_name string, area int) *coopMember_harvest {
	m := &coopMember_harvest{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	//infor
	m.area = area

	//state
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("CHECK", 0)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("check")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("player")

	return m
}

func (m *coopMember_harvest) Ext_trans(port string, msg *system.SysMessage) {
	//player가 해당 셀에 왔음
	if port == "check" {

		m.executor.Cancel_rescheduling()

		m.executor.Cur_state = "CHECK"
	}

}

func (m *coopMember_harvest) Output() *system.SysMessage {

	return nil
}

func (m *coopMember_harvest) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "CHECK" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IDLE"
	}
}

//Shipment
type coopMember_ship struct {
	executor *executor.BehaviorModelExecutor
	msg      *system.SysMessage
}

func AM_ship(instance_time, destruct_time float64, name, engine_name string) *coopMember_ship {
	m := &coopMember_ship{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	//infor

	//statef
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("IN", 0)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("in")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("check")

	return m
}

func (m *coopMember_ship) Ext_trans(port string, msg *system.SysMessage) {
	//player가 해당 셀에 왔음
	if port == "in" {
		m.executor.Cancel_rescheduling()

		m.executor.Cur_state = "IN"
	}
}

func (m *coopMember_ship) Output() *system.SysMessage {
	//check 에게 출력을 보내서 동작시킴
	fmt.Println("State: cell in")
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "check")

	return msg
}

func (m *coopMember_ship) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "IN" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IDLE"
	}
}
