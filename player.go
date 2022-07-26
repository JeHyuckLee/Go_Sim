package main

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
	"evsim_golang/system"
	"fmt"
)

type cm_player struct {
	am_move  *move
	am_think *think
}

func create_player(instance_time, destruct_time float64, name, engine_name string, ix, iy int) *cm_player {
	player := cm_player{}
	player.am_move = AM_move(instance_time, destruct_time, name, engine_name)
	player.am_think = AM_think(instance_time, destruct_time, name, engine_name)

	return &player
}

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
		m.msg_list = append(m.msg_list, 0)
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
	ahead    Ahead
	pos      pos
	cell_msg
	nx, ny int
}

func AM_think(instance_time, destruct_time float64, name, engine_name string) *think {
	m := think{}
	m.pos.x = 0
	m.pos.y = 0
	m.set_Ahead("south")
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

func (m *think) Ext_trans(port string, msg *system.SysMessage) {
	//cell에게 입력을 받은 정보를 토대로 어디로 이동할지 생각
	if port == "player" {
		// cell 입력받은 정보를 저장
		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		m.cell_msg = data[0].(cell_msg)
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

func (m *think) Int_trans() {
	if m.executor.Cur_state == "THINK" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "THINK"
	}
}

func (m *think) set_Ahead(ahead string) {
	switch ahead {
	case "north":
		m.ahead.front = "north"
		m.ahead.back = "south"
		m.ahead.left = "east"
		m.ahead.right = "west"
	case "south":
		m.ahead.front = "south"
		m.ahead.back = "north"
		m.ahead.left = "west"
		m.ahead.right = "east"
	case "east":
		m.ahead.front = "east"
		m.ahead.back = "west"
		m.ahead.left = "south"
		m.ahead.right = "north"
	case "west":
		m.ahead.front = "west"
		m.ahead.back = "east"
		m.ahead.left = "north"
		m.ahead.right = "south"
	}
}

func (m *think) turnLeft() {
	switch m.ahead.front {
	case "north":
		m.set_Ahead("east")
	case "south":
		m.set_Ahead("west")
	case "east":
		m.set_Ahead("south")
	case "west":
		m.set_Ahead("north")
	}
}

func (m *think) turnRight() {
	switch m.ahead.front {
	case "north":
		m.set_Ahead("west")
	case "south":
		m.set_Ahead("east")
	case "east":
		m.set_Ahead("north")
	case "west":
		m.set_Ahead("south")
	}
}

func (m *think) set_position(x int, y int) {
	m.pos.x = x
	m.pos.y = y
}

func (m *think) get_position() pos {
	return m.pos
}
