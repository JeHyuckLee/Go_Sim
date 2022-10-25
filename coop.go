package main

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
	"evsim_golang/system"
	"fmt"
)

//Cooperative
type cm_coop struct {
	am_seed    *coop_ware
	am_harvest *coop_management
	am_ship    *coop_shipment
}

func CM_coop(instance_time, destruct_time float64, name, engine_name string, area int) *cm_coop {

	cell := cm_coop{}
	cell.am_seed = AM_ware(instance_time, destruct_time, name, engine_name, area)
	cell.am_harvest = AM_management(instance_time, destruct_time, name, engine_name, area)
	cell.am_ship = AM_shipment(instance_time, destruct_time, name, engine_name)

	return &cell
}

//Warehousing
type coop_ware struct {
	executor *executor.BehaviorModelExecutor
	area     int
	harvest  int
	req_item interface{}
	msg      *system.SysMessage
}

func AM_ware(instance_time, destruct_time float64, name, engine_name string, area int) *coop_ware {
	m := &coop_ware{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	//infor
	m.area = area
	m.harvest = 3

	//statef
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("WARE", 0)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("in")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("manage")

	return m
}

func (m *coop_ware) Ext_trans(port string, msg *system.SysMessage) {
	//파종이 필요하다고 요청이 옴
	if port == "in" {
		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		m.req_item = data[0]

		m.executor.Cur_state = "WARE"
	}
}

func (m *coop_ware) Output() *system.SysMessage {
	//가능한 수확량선에서 필요한 만큼 파종을 함
	fmt.Println("State: Warehousing in")
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "manage")

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

func (m *coop_ware) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "SEEDING" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IDLE"
	}
}

//Management
type coop_management struct {
	executor *executor.BehaviorModelExecutor
	area     int
	harvest  int
	msg      *system.SysMessage
}

func AM_management(instance_time, destruct_time float64, name, engine_name string, area int) *coop_management {
	m := &coop_management{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	//infor
	m.area = area

	//state
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("HARVEST", 0)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("harvest")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("shipment")

	return m
}

func (m *coop_management) Ext_trans(port string, msg *system.SysMessage) {
	//player가 해당 셀에 왔음
	if port == "harvest" {
		fmt.Println("[Seeding] => [Harvest]")
		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		m.harvest = data[0].(int)

		m.executor.Cur_state = "HARVEST"
	}

}

func (m *coop_management) Output() *system.SysMessage {
	fmt.Println("Harvest...")
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "harvest")
	msg.Insert(m.harvest)

	return msg
}

func (m *coop_management) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "HARVEST" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IDLE"
	}
}

//Shipment
type coop_shipment struct {
	executor *executor.BehaviorModelExecutor
	shipment int
	msg      *system.SysMessage
}

func AM_shipment(instance_time, destruct_time float64, name, engine_name string) *coop_shipment {
	m := &coop_shipment{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	//infor

	//statef
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("SHIPMENT", 0)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("shipment")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("in")

	return m
}

func (m *coop_shipment) Ext_trans(port string, msg *system.SysMessage) {
	//player가 해당 셀에 왔음
	if port == "shipment" {
		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		m.shipment = data[0].(int)

		m.executor.Cur_state = "SHIPMENT"
	}
}

func (m *coop_shipment) Output() *system.SysMessage {
	//check 에게 출력을 보내서 동작시킴
	fmt.Println("Shipment...")
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "in")
	msg.Insert(m.shipment)

	return msg
}

func (m *coop_shipment) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "SHIPMENT" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IDLE"
	}
}
