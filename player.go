package main

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
	"evsim_golang/system"
	"fmt"
)

//player 의 원자모델 move
type move struct {
	executor    *executor.BehaviorModelExecutor
	msg_list    []interface{}
	current_pos pos
	next_pos    pos
}

//atomic model
func AM_move(instance_time, destruct_time float64, name, engine_name string) *move {
	m := move{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = &m

	//state
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("MOVE", 1)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("start")
	m.executor.Behaviormodel.CoreModel.Insert_input_port("think")

	return &m
}

func (m *move) set_position(x int, y int) {
	m.current_pos.x = x
	m.current_pos.y = y
}

func (m *move) get_position() pos {
	return m.current_pos
}

func (m *move) insert_Player_Output_Port(port_name string) {
	m.executor.Behaviormodel.CoreModel.Insert_output_port(port_name)
}

func (m *move) Int_trans() {
	if m.executor.Cur_state == "MOVE" {
		//이동
		m.set_position(m.next_pos.x, m.next_pos.y)

		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "MOVE"
	}
}

func (m *move) Ext_trans(port string, msg *system.SysMessage) {
	//think로 부터 입력받아 해당하는 cell로 이동
	if port == "start" {
		m.set_position(0, 0)
		m.executor.Cur_state = "MOVE"
	}

	if port == "think" {

		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		p := data[0].(pos)

		//다음움직일곳
		m.next_pos.x = p.x
		m.next_pos.y = p.y
		m.executor.Cur_state = "MOVE"
	}
}

func (m *move) Output() *system.SysMessage {
	//그 해당하는 cell로 이동 해당 셀에 입력을 보냄
	output_port := fmt.Sprintf("{%d,%d}", m.get_position().x, m.get_position().y)
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), output_port)
	msg.Insert(m.msg_list[0])
	return msg
}

//player의 원자모델
type think struct {
	executor *executor.BehaviorModelExecutor
	msg_list []interface{}
	pos      pos
	nx, ny   int
}

func AM_think(instance_time, destruct_time float64, name, engine_name string) *think {
	m := think{}
	m.pos.x = 0
	m.pos.y = 0
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = &m

	//state
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("THINK", 1)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("player")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("move")

	return &m
}

func (m *think) Int_trans() {
	if m.executor.Cur_state == "THINK" {
		//cell로 부터 갈수있는 위치와 방향을 받아서 어디로갈지 정하는 로직
		//m.set_position(nx,ny)
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "THINK"
	}
}

func (m *think) Ext_trans(port string, msg *system.SysMessage) {
	//cell에게 입력을 받은 정보를 토대로 어디로 이동할지 생각
	if port == "player" {
		// cell 입력받은 정보를 저장
		m.executor.Cur_state = "THINK"
	}
}

func (m *think) Output() *system.SysMessage {
	//이동할 위치를 전송
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "move")
	m.msg_list = append(m.msg_list, m.get_position())
	msg.Insert(m.msg_list[0])
	return msg
}

func (m *think) set_position(x int, y int) {
	m.pos.x = x
	m.pos.y = y
}

func (m *think) get_position() pos {
	return m.pos
}
