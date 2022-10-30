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
	buy_log  []int
	msg      *system.SysMessage
}

func AM_buyer(instance_time, destruct_time float64, name, engine_name string, buy int) *buyer {
	m := &buyer{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m
	m.buy = buy
	for i := 0; i < 12; i++ {
		m.buy_log = append(m.buy_log, 0)
	}

	m.executor.Behaviormodel.Insert_state("IDLE", 120)
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

	Sales_date := m.executor.Get_req_time()

	date_cal(m.buy, int(Sales_date), m.buy_log)

	//
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

func date_cal(amount int, date int, a []int) *[]int {
	if date >= 0 && date < 31 {
		a[0] += amount
	} else if date < 58 {
		a[1] += amount
	} else if date < 89 {
		a[2] += amount
	} else if date < 119 {
		a[3] += amount
	} else if date < 150 {
		a[4] += amount
	} else if date < 180 {
		a[5] += amount
	} else if date < 211 {
		a[6] += amount
	} else if date < 242 {
		a[7] += amount
	} else if date < 272 {
		a[8] += amount
	} else if date < 303 {
		a[9] += amount
	} else if date < 333 {
		a[10] += amount
	} else if date < 364 {
		a[11] += amount
	}

	return &a
}
