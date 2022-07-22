package main

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
	"evsim_golang/system"
)

//cell의 원자모델
type cellout struct {
	executor *executor.BehaviorModelExecutor
	msg_list []interface{}
}

func (m *cellout) Int_trans() {
	//상태변화
}

func (m *cellout) Ext_trans(port string, msg *system.SysMessage) {
	//check 에게 정보를 받음
}

func (m *cellout) Output() *system.SysMessage {
	//player 에게 전송
}

func AM_cellout() *cellout {
	m := cellout{}
	m.executor = executor.NewExecutor(0, definition.Infinite, "out", "gosim")
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
type cellin struct {
	executor *executor.BehaviorModelExecutor
	msg_list []interface{}
}

func (m *cellin) Int_trans() {
	//상태변화
}

func (m *cellin) Ext_trans(port string, msg *system.SysMessage) {
	//player가 해당 셀에 왔음
}

func (m *cellin) Output() *system.SysMessage {
	//check 에게 출력을 보내서 동작시킴
}

func AM_cellin() *cellin {
	m := cellin{}
	m.executor = executor.NewExecutor(0, definition.Infinite, "in", "gosim")
	m.executor.AbstractModel = &m

	//state
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("CHECK", 1)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("cell")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("check")

	return &m
}

//cell의 원자모델
type check struct {
	executor *executor.BehaviorModelExecutor
	msg_list []interface{}
	block    bool
}

func (m *check) set_block(b bool) {
	m.block = b
}
func (m *check) get_block() bool {
	return m.block
}

func (m *check) Int_trans() {
	//상태천이
}

func (m *check) Ext_trans(port string, msg *system.SysMessage) {
	//in에게 입력을 받으면 check 상태로 가고 연결된 애들에게 입력을 보냄
	//NEWS 포트로 입력을 받으면 out 상태로 가고 OUT에게 입력을 보냄
}

func (m *check) Output() *system.SysMessage {
	//in에게 입력을 받으면 NEWS 포트중 연결된 포트로 출력
	//NEWS포트 로 입력이 들어오면 입력된 정보를 OuT 에게 전송
}

func AM_check() *check {
	//맵 모델
	m := check{}
	m.executor = executor.NewExecutor(0, definition.Infinite, "maze", "gosim")
	m.executor.AbstractModel = &m
	m.block = false

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
