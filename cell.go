package main

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
	"evsim_golang/system"
	"fmt"
)

//cell의 원자모델
type cellOut struct {
	executor *executor.BehaviorModelExecutor
	msg_list []interface{}
}

func (m *cellOut) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "OUT" && len(m.msg_list) == 0 {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "OUT"
	}
}

func (m *cellOut) Ext_trans(port string, msg *system.SysMessage) {
	//check 에게 정보를 받음
	if port == "check" {

	}
}

func (m *cellOut) Output() *system.SysMessage {
	//player 에게 전송
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "player")

	return msg
}

func AM_cellOut(instance_time, destruct_time float64, name, engine_name string) *cellOut {
	m := cellOut{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = &m

	//state
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("OUT", 0)
	m.executor.Init_state("IDLE")
	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("check")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("player")
	return &m
}

//cell의 원자모델
type cellIn struct {
	executor    *executor.BehaviorModelExecutor
	player_list []interface{}
}

func (m *cellIn) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "IN" && len(m.player_list) == 0 {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IN"
	}
}

func (m *cellIn) Ext_trans(port string, msg *system.SysMessage) {
	//player가 해당 셀에 왔음
	if port == m.executor.Behaviormodel.CoreModel.Get_name() {
		fmt.Println("Cell: ", m.executor.Behaviormodel.CoreModel.Get_name())
		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		m.player_list = append(m.player_list, data...)
		m.executor.Cur_state = "IN"
	}

}

func (m *cellIn) Output() *system.SysMessage {
	//check 에게 출력을 보내서 동작시킴
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "check")
	msg.Insert(m.player_list[0])
	m.player_list = remove(m.player_list, 0)

	return msg
}

func AM_cellIn(instance_time, destruct_time float64, name, engine_name string) *cellIn {
	m := cellIn{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = &m

	//state
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("IN", 1)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port(m.executor.Behaviormodel.CoreModel.Get_name())
	m.executor.Behaviormodel.CoreModel.Insert_output_port("check")

	return &m
}

//cell의 원자모델
type check struct {
	executor   *executor.BehaviorModelExecutor
	block_list []interface{}
	block      bool
	x          int
	y          int
	count      int
	Nflag      bool
	Sflag      bool
	Eflag      bool
	Wflag      bool
}

func (m *check) set_position(x int, y int) {
	m.x = x
	m.y = y
}

func (m *check) get_position() (int, int) {
	return m.x, m.y
}

func (m *check) set_block(b bool) {
	m.block = b
}
func (m *check) get_block() bool {
	return m.block
}

func (m *check) Int_trans() {
	//상태천이
	if m.executor.Cur_state == "CHECK" && m.count == 5 {
		m.executor.Cur_state = "IDLE"
	} else if m.executor.Cur_state == "SET" {
		m.executor.Cur_state = "IDLE"
	} else if m.executor.Cur_state == "OUT" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "CHECK"
	}
}

func (m *check) Ext_trans(port string, msg *system.SysMessage) {
	//in에게 입력을 받으면 check 상태로 가고 연결된 애들에게 입력을 보냄
	//NEWS 포트로 입력을 받으면 SET 상태로 가고 각 셀 정보 저장
	//전부 탐색시 OUT 상태로 가고 out에 정보 전달
	if port == "in" {
		m.executor.Cur_state = "CHECK"
	} else if m.count == 4 {
		m.executor.Cur_state = "OUT"
		m.count++
	} else {
		m.executor.Cur_state = "SET"
		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		m.block_list = append(m.block_list, data...)
	}
}

func (m *check) Output() *system.SysMessage {
	//in에게 입력을 받으면 NEWS 포트중 연결된 포트로 출력
	//NEWS포트 로 입력이 들어오면 입력된 정보를 OUT 에게 전송
	if m.executor.Cur_state == "CHECK" {
		if m.Nflag == false {
			msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "north")
			m.Nflag = true
			m.count++
			return msg
		} else if m.Sflag == false {
			msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "south")
			m.Sflag = true
			m.count++
			return msg
		} else if m.Eflag == false {
			msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "east")
			m.Eflag = true
			m.count++
			return msg
		} else if m.Wflag == false {
			msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "west")
			m.Wflag = true
			m.count++
			return msg
		}
		return nil
	} else if m.executor.Cur_state == "SET" {
		return nil
	}

	// Cur_state = OUT
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "out")
	msg.Insert(m.block_list)
	m.block_list = nil

	return msg
}

func AM_check(instance_time, destruct_time float64, name, engine_name string, px int, py int) *check {
	//맵 모델
	m := check{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = &m
	m.block = false
	m.x = px
	m.y = py
	m.Nflag = false
	m.Sflag = false
	m.Eflag = false
	m.Wflag = false

	//state

	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	//OUT으로부터 입력이 오면 IDLE -> CHECK
	m.executor.Behaviormodel.Insert_state("CHECK", 0)
	//Chcek한 셀 정보 넘기기
	m.executor.Behaviormodel.Insert_state("OUT", 0)
	//NEWS로부터 입력이 오면 check한 셀 정보 저장
	m.executor.Behaviormodel.Insert_state("SET", 0)

	//input port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("north")
	m.executor.Behaviormodel.CoreModel.Insert_input_port("south")
	m.executor.Behaviormodel.CoreModel.Insert_input_port("east")
	m.executor.Behaviormodel.CoreModel.Insert_input_port("west")
	m.executor.Behaviormodel.CoreModel.Insert_input_port("in")
	//output port
	m.executor.Behaviormodel.CoreModel.Insert_output_port("north")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("south")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("east")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("west")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("out")

	return &m
}

func remove(slice []interface{}, s int) []interface{} {
	return append(slice[:s], slice[s+1:]...)
}
