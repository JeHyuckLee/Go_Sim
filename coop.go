package main

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
	"evsim_golang/system"
	"fmt"
)

// Cooperative
type cm_coop struct {
	inventory     []tomato
	am_ware       *coop_ware
	am_management *coop_management
	am_shipment   *coop_shipment
}

func CM_coop(instance_time, destruct_time float64, name, engine_name string, storage_period int) *cm_coop {

	coop := cm_coop{}
	// file, err := os.Create("./output.csv")
	// if err != nil {
	// 	panic(err)
	// }

	coop.am_ware = AM_ware(instance_time, destruct_time, name, engine_name, &coop.inventory)
	coop.am_management = AM_management(instance_time, destruct_time, name, engine_name, storage_period, &coop.inventory)
	coop.am_shipment = AM_shipment(instance_time, destruct_time, name, engine_name, &coop.inventory)

	return &coop
}

// Warehousing
type coop_ware struct {
	executor  *executor.BehaviorModelExecutor
	inventory *[]tomato
	received  tomato
	msg       *system.SysMessage
}

func AM_ware(instance_time, destruct_time float64, name, engine_name string, inventory *[]tomato) *coop_ware {
	m := &coop_ware{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	m.inventory = inventory
	//statef
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("WARE", 1)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("in")

	return m
}

func (m *coop_ware) Ext_trans(port string, msg *system.SysMessage) {
	//파종이 필요하다고 요청이 옴
	if port == "in" {
		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		m.received = data[0].(tomato)
		*m.inventory = append(*m.inventory, m.received)
		fmt.Println("[Warehousing] Current inventory : ", total_tomato(m.inventory))
		m.executor.Cur_state = "WARE"
	}
}

func (m *coop_ware) Output() *system.SysMessage {

	fmt.Println("State: Warehousing in")

	return nil
}

func (m *coop_ware) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "WARE" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IDLE"
	}
}

// Management
type coop_management struct {
	executor       *executor.BehaviorModelExecutor
	storage_period int
	inventory      *[]tomato
	msg            *system.SysMessage
}

func AM_management(instance_time, destruct_time float64, name, engine_name string, storage_period int, inventory *[]tomato) *coop_management {
	m := &coop_management{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	m.inventory = inventory
	//infor
	m.storage_period = storage_period

	//state
	m.executor.Behaviormodel.Insert_state("IDLE", 1)
	m.executor.Behaviormodel.Insert_state("CHECK", 0)
	m.executor.Init_state("IDLE")

	return m
}

func (m *coop_management) Ext_trans(port string, msg *system.SysMessage) {
}

func (m *coop_management) Output() *system.SysMessage {
	return nil
}

func (m *coop_management) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "IDLE" {
		if len((*m.inventory)) > 0 {
			// Sort_tomato(m.inventory)

			if (*m.inventory)[0].Period <= 0 {
				for k, v := range *m.inventory {
					if v.Period <= 0 {
						(*m.inventory) = remove_tomato(m.inventory, k)
						fmt.Println("보관기간 지나서 버림 : ", v.Quantity)
					}
				}

			} else {
				for k, _ := range *m.inventory {
					(*m.inventory)[k].Next_day()
				}

			}
		}
		m.executor.Cur_state = "CHECK"
	} else {
		m.executor.Cur_state = "IDLE"
	}
}

// Shipment
type coop_shipment struct {
	executor         *executor.BehaviorModelExecutor
	shipment_qantity int
	inventory        *[]tomato
	msg              *system.SysMessage
}

func AM_shipment(instance_time, destruct_time float64, name, engine_name string, inventory *[]tomato) *coop_shipment {
	m := &coop_shipment{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	//infor
	m.inventory = inventory
	//statef
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("SHIPMENT", 1)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("shipment")

	return m
}

func (m *coop_shipment) Ext_trans(port string, msg *system.SysMessage) {
	//
	if port == "shipment" {
		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		if len((*m.inventory)) > 0 {
			m.shipment_qantity = data[0].(int)

			fmt.Println("Shipment Quantity : ", m.shipment_qantity)

			Sales(m.shipment_qantity, m.inventory)
			fmt.Println("[Sales] Current Inventory : ", total_tomato(m.inventory))
			m.executor.Cur_state = "SHIPMENT"
		} else {
			fmt.Println("재고가없어서 못팜")

		}

	}
}

func (m *coop_shipment) Output() *system.SysMessage {
	//check 에게 출력을 보내서 동작시킴
	return nil
}

func (m *coop_shipment) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "SHIPMENT" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IDLE"
	}
}
