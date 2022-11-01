package main

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
	"evsim_golang/system"
	"fmt"
)

// Seeding
type buyer struct {
	executor *executor.BehaviorModelExecutor
	buyer    buyer_infor
	buy      int
	cnt      int
	term     int
	i        int
	msg      *system.SysMessage
}

func AM_buyer(instance_time, destruct_time float64, name, engine_name string, buy buyer_infor) *buyer {
	m := &buyer{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m
	m.i = 0
	m.buyer = buy
	m.cnt = buy.cnt
	m.term = 100 / buy.cnt

	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)

	m.executor.Behaviormodel.Insert_state("BUY", 1) //나중에 멤버에게 입력받아서 집어넣어야함
	m.executor.Behaviormodel.Insert_state("REQ", float64(m.term))
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("start")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("buy")

	return m
}

func (m *buyer) Ext_trans(port string, msg *system.SysMessage) {

	if port == "start" {
		m.executor.Cancel_rescheduling()
		fmt.Println("buyer Start")
		m.executor.Cur_state = "BUY"
	}
}

func (m *buyer) Output() *system.SysMessage {
	//가능한 수확량선에서 필요한 만큼 파종을 함
	if m.executor.Cur_state == "BUY" {
		m.i++
		m.buy = int(rand_crop(m.buyer.aver, m.buyer.std))
		fmt.Println("Buyer: buy a tomato : ", m.buy)

		msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "buy")
		msg.Insert(m.buy)
		return msg
	} else {
		return nil
	}

}

func (m *buyer) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "BUY" {
		m.executor.Cur_state = "REQ"
	} else if m.executor.Cur_state == "REQ" && m.i < m.cnt {
		m.executor.Cur_state = "BUY"
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

func date_month(date int) int {
	if date >= 0 && date < 31 {
		return 1
	} else if date < 58 {
		return 2
	} else if date < 89 {
		return 3
	} else if date < 119 {
		return 4
	} else if date < 150 {
		return 5
	} else if date < 180 {
		return 6
	} else if date < 211 {
		return 7
	} else if date < 242 {
		return 8
	} else if date < 272 {
		return 9
	} else if date < 303 {
		return 10
	} else if date < 333 {
		return 11
	} else if date < 364 {
		return 12
	} else {
		return 0
	}
}
