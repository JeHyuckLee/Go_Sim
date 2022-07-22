package main

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
	"evsim_golang/system"
)

type Dir int

const (
	Dir_UP = itoa
	Dir_LEFT
	Dir_DOWN
	Dir_RIGHT
	DIR_COUNT
)

//player 의 원자모델 move
type move struct {
	executor *executor.BehaviorModelExecutor
	msg_list []interface{}
}

func (m *move) Int_trans() {

}

func (m *move) Ext_trans(port string, msg *system.SysMessage) {
	//think 로 부터 입력받아 해당하는 cell로 이동
}

func (m *move) Output() *system.SysMessage {
	//그 해당하는 cell에게 입력을 보냄
}

//atomic model
func AC_move() *move {
	m := move{}
	m.executor = executor.NewExecutor(0, definition.Infinite, "move", "gosim")
	m.executor.AbstractModel = &m

	//state
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("MOVE", 1)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("start")
	m.executor.Behaviormodel.CoreModel.Insert_input_port("think")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("cell")

	return &m
}

//player의 원자모델
type think struct {
	executor *executor.BehaviorModelExecutor
	msg_list []interface{}
}

func (m *think) Int_trans() {

}

func (m *think) Ext_trans(port string, msg *system.SysMessage) {
	//cell에게 입력을 받은 정보를 토대로 어디로 이동할지 생각
}

func (m *think) Output() *system.SysMessage {
	//이동할 위치를 전송
}

func AC_think() *think {
	m := think{}
	m.executor = executor.NewExecutor(0, definition.Infinite, "think", "gosim")
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
