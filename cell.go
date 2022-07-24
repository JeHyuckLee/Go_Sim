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
	executor *executor.BehaviorModelExecutor
	msg_list []interface{}
	block    bool
	x        int
	y        int
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
	if m.executor.Cur_state == "CHECK" && len(m.msg_list) == 0 {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "CHECK"
	}
}

func (m *check) Ext_trans(port string, msg *system.SysMessage) {
	//in에게 입력을 받으면 check 상태로 가고 연결된 애들에게 입력을 보냄
	//NEWS 포트로 입력을 받으면 out 상태로 가고 OUT에게 입력을 보냄
	if port == "in" {
		m.executor.Cur_state = "CHECK"
	} else {

	}
}

func (m *check) Output() *system.SysMessage {
	//in에게 입력을 받으면 NEWS 포트중 연결된 포트로 출력
	//NEWS포트 로 입력이 들어오면 입력된 정보를 OuT 에게 전송
	if m.executor.Cur_state == "CHECK" {

	}
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "in")
	if m.executor.Cur_state == "OUT" {
		if m.block == true {
			fmt.Println("cell[%d][%d]", m.x, m.y)
		} else {
			m.executor.Cur_state = "OUT"
		}
	}
	msg.Insert(m.msg_list[0])
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

	//state

	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	//OUT으로부터 입력이 오면 IDLE -> CHECK
	m.executor.Behaviormodel.Insert_state("CHECK", 0)
	//NEWS로부터 입력이 오면 IDLE->OUT
	m.executor.Behaviormodel.Insert_state("OUT", 0)

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
