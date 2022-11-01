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

	coop.am_ware = AM_ware(instance_time, destruct_time, name, engine_name, &coop.inventory)
	coop.am_management = AM_management(instance_time, destruct_time, name, engine_name, storage_period, &coop.inventory)
	coop.am_shipment = AM_shipment(instance_time, destruct_time, name, engine_name, &coop.inventory)

	return &coop
}

// Warehousing
type coop_ware struct {
	executor  *executor.BehaviorModelExecutor
	inventory *[]tomato

	ware []int

	received tomato
	msg      *system.SysMessage
}

func AM_ware(instance_time, destruct_time float64, name, engine_name string, inventory *[]tomato) *coop_ware {
	m := &coop_ware{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	m.inventory = inventory

	for i := 0; i < 12; i++ {
		m.ware = append(m.ware, 0)
	}
	//
	db := GetConnector()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 12; i++ {
		results, err := db.Exec("UPDATE Simulate_Sales SET Warehousing_amount = ? WHERE Sales_date = ?", m.ware[i], i+1)
		if err != nil {
			panic(err.Error())
		}
		n, err := results.RowsAffected()
		if n == 1 {
		}
	}

	//statef
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("WARE", 0)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("in")

	return m
}

func (m *coop_ware) Ext_trans(port string, msg *system.SysMessage) {

	if port == "in" {
		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		m.received = data[0].(tomato)
		*m.inventory = append(*m.inventory, m.received)

		fmt.Println("[coop Warehousing] Current inventory : ", total_tomato(m.inventory))

		m.executor.Cur_state = "WARE"
	}
}

func (m *coop_ware) Output() *system.SysMessage {

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

	trash []int

	msg *system.SysMessage
}

func AM_management(instance_time, destruct_time float64, name, engine_name string, storage_period int, inventory *[]tomato) *coop_management {
	m := &coop_management{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	m.inventory = inventory
	//infor
	m.storage_period = storage_period

	for i := 0; i < 12; i++ {
		m.trash = append(m.trash, 0)
	}

	db := GetConnector()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 12; i++ {
		results, err := db.Exec("UPDATE Simulate_Sales SET trash_amount = ? WHERE Sales_date = ?", m.trash[i], i+1)
		if err != nil {
			panic(err.Error())
		}
		n, err := results.RowsAffected()
		if n == 1 {
		}
	}

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
					if v.Period <= 0 && v.Quantity != 0 {
						(*m.inventory) = remove_tomato(m.inventory, k)

						fmt.Println("[coop] 보관기간 지나서 버려짐 : ", v.Quantity)
						fmt.Println("[coop] 남은 토마토 : ", total_tomato(m.inventory))
						Sales_date := m.executor.Get_req_time()
						month := date_month(int(Sales_date))
						//
						db := GetConnector()
						defer db.Close()

						err := db.Ping()
						if err != nil {
							panic(err)
						}

						results, err := db.Exec("UPDATE Simulate_Sales SET trash_amount = ? WHERE Sales_date = ?", v.Quantity, month)
						if err != nil {
							panic(err.Error())
						}
						n, err := results.RowsAffected()
						if n == 1 {
						}
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
	buy_log          []int
	msg              *system.SysMessage
}

func AM_shipment(instance_time, destruct_time float64, name, engine_name string, inventory *[]tomato) *coop_shipment {
	m := &coop_shipment{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	//infor
	m.inventory = inventory
	for i := 0; i < 12; i++ {
		m.buy_log = append(m.buy_log, 0)
	}
	db := GetConnector()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		panic(err)
	}
	for i := 0; i < 12; i++ {
		results, err := db.Exec("UPDATE Simulate_Sales SET Sales_amount = ? WHERE Sales_date = ?", m.buy_log[i], i+1)
		if err != nil {
			panic(err.Error())
		}
		n, err := results.RowsAffected()
		if n == 1 {
		}
	}
	//statef
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("SHIPMENT", 0)
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
		m.shipment_qantity = data[0].(int)
		m.executor.Cur_state = "SHIPMENT"
	}
}

func (m *coop_shipment) Output() *system.SysMessage {
	//check 에게 출력을 보내서 동작시킴
	total := total_tomato(m.inventory)
	if len((*m.inventory)) > 0 {

		Sales_date := m.executor.Get_req_time()
		month := date_month(int(Sales_date))
		fmt.Println("salse date", Sales_date)
		if m.shipment_qantity > total {
			oversell := m.shipment_qantity - total

			db := GetConnector()
			defer db.Close()

			err := db.Ping()
			if err != nil {
				panic(err)
			}

			results, err := db.Exec("UPDATE Simulate_Sales SET Oversell_request = Oversell_request + ? WHERE Sales_date = ?", oversell, month)
			if err != nil {
				panic(err.Error())
			}
			n, err := results.RowsAffected()
			if n == 1 {
			}

			fmt.Println("[coop] ShipmentQuantity : ", total)
			fmt.Println(total)

			//

			result, err := db.Exec("UPDATE Simulate_Sales SET Sales_amount = Sales_amount + ? WHERE Sales_date = ?", total, month)
			if err != nil {
				panic(err.Error())
			}
			a, err := result.RowsAffected()
			if a == 1 {
			}

			Sales(total_tomato(m.inventory), m.inventory)
		} else {
			fmt.Println("[coop] Shipment Quantity : ", m.shipment_qantity)
			fmt.Println(m.shipment_qantity)
			db := GetConnector()
			defer db.Close()

			err := db.Ping()
			if err != nil {
				panic(err)
			}

			results, err := db.Exec("UPDATE Simulate_Sales SET Sales_amount = Sales_amount + ? WHERE Sales_date = ?", m.shipment_qantity, month)
			if err != nil {
				panic(err.Error())
			}
			n, err := results.RowsAffected()
			if n == 1 {
			}
			Sales(m.shipment_qantity, m.inventory)
		}
		fmt.Println("[coop Sales] Current Inventory : ", total_tomato(m.inventory))
		m.executor.Cur_state = "SHIPMENT"
	}

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
