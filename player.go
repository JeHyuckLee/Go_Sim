package main

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
	"evsim_golang/system"
	"fmt"
)

type Dir int

// const (
// 	Dir_UP = itoa
// 	Dir_LEFT
// 	Dir_DOWN
// 	Dir_RIGHT
// 	DIR_COUNT
// )

//player 의 원자모델 move
type move struct {
	executor *executor.BehaviorModelExecutor
	msg_list []interface{}
	x        int
	y        int
	portName string
}

func (m *move) set_position(x int, y int) {
	m.x = x
	m.y = y
}

func (m *move) get_position() (int, int) {
	return m.x, m.y
}

//atomic model
func AM_move(instance_time, destruct_time float64, name, engine_name string) *move {
	m := move{}
	m.portName = fmt.Sprintf("{%n,%n}", m.x, m.y)
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = &m

	//state
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("MOVE", 1)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("start")
	m.executor.Behaviormodel.CoreModel.Insert_input_port("think")
	m.executor.Behaviormodel.CoreModel.Insert_output_port(m.portName)

	return &m
}

func (m *move) Int_trans() {
	if m.executor.Cur_state == "MOVE" && len(m.msg_list) == 0 {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "MOVE"
	}
}

func (m *move) Ext_trans(port string, msg *system.SysMessage) {
	//think로 부터 입력받아 해당하는 cell로 이동
	if port == "start" {
		m.executor.Cur_state = "MOVE"
		m.set_position(0, 0)
	}

	if port == "think" {
		m.executor.Cur_state = "THINK"
		m.set_position(0, 0)
	}
}

func (m *move) Output() *system.SysMessage {
	//그 해당하는 cell로 이동 해당 셀에 입력을 보냄

	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), m.portName)
	msg.Insert(m.msg_list[0])
	return msg
}

//player의 원자모델
type think struct {
	executor *executor.BehaviorModelExecutor
	msg_list []interface{}
	x        int
	y        int
}

func AM_think(instance_time, destruct_time float64, name, engine_name string) *think {
	m := think{}
	m.x = 0
	m.y = 0
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
	if m.executor.Cur_state == "THINK" && len(m.msg_list) == 0 {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "THINK"
	}
}

func (m *think) Ext_trans(port string, msg *system.SysMessage) {
	//cell에게 입력을 받은 정보를 토대로 어디로 이동할지 생각
	if port == "player" {
		m.executor.Cur_state = "THINK"
	}
}

func (m *think) Output() *system.SysMessage {
	//이동할 위치를 전송
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "move")
	msg.Insert(m.msg_list[0])
	return msg
}

func (m *think) set_position(x int, y int) {
	m.x = x
	m.y = y
}

func (m *think) get_position() (int, int) {
	return m.x, m.y
}
