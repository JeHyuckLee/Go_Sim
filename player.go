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
	player.am_move = AM_move(instance_time, destruct_time, name, engine_name, ix, iy)
	player.am_think = AM_think(instance_time, destruct_time, name, engine_name)

	return &player
}

//player 의 원자모델 move
type move struct {
	executor    *executor.BehaviorModelExecutor
	ahead       Ahead
	current_pos pos
}

//atomic model
func AM_move(instance_time, destruct_time float64, name, engine_name string, ix, iy int) *move {
	m := move{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = &m
	m.current_pos.x = ix
	m.current_pos.y = iy

	//state
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("MOVE", 1)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("start")
	m.executor.Behaviormodel.CoreModel.Insert_input_port("think")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("in")

	return &m
}

func (m *move) move_player(dir Dir) {
	switch dir {
	case Dir(0):
		m.set_position(m.current_pos.x, m.current_pos.y-1)
	case Dir(1):
		m.set_position(m.current_pos.x+1, m.current_pos.y)
	case Dir(2):
		m.set_position(m.current_pos.x-1, m.current_pos.y)
	case Dir(3):
		m.set_position(m.current_pos.x, m.current_pos.y+1)
	}
}

func (m *move) set_position(x int, y int) {
	m.current_pos.x = x
	m.current_pos.y = y
}

func (m *move) get_position() pos {
	return m.current_pos
}

func (m *move) Ext_trans(port string, msg *system.SysMessage) {
	//think로 부터 입력받아 해당하는 cell로 이동
	if port == "start" {
		fmt.Println("Hi Maze")
		m.executor.Cur_state = "MOVE"
	}

	if port == "think" {
		fmt.Println("State: player think -> move")
		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		m.ahead = data[0].(Ahead)
		if m.current_pos.x == 98 && m.current_pos.y == 98 {
			fmt.Println("Arrive Destination!!!")
			m.executor.Cur_state = "IDLE"
		} else {
			m.move_player(m.ahead.front)
			m.executor.Cur_state = "MOVE"
		}
		// 플레이어 이동
	}
}

func (m *move) Output() *system.SysMessage {
	//그 해당하는 cell로 이동 해당 셀에 입력을 보냄
	output_port := fmt.Sprintf("{%d,%d}", m.get_position().x, m.get_position().y)

	fmt.Println(output_port)
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "in")
	msg.Insert(m.current_pos)
	return msg
}

func (m *move) Int_trans() {
	if m.executor.Cur_state == "MOVE" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "MOVE"
	}
}

//player의 원자모델
type think struct {
	executor   *executor.BehaviorModelExecutor
	ahead      Ahead
	pos        pos
	input_msg  []cell_msg
	nx, ny     int
	flag       bool
	right_flag bool
}

func AM_think(instance_time, destruct_time float64, name, engine_name string) *think {
	m := think{}
	m.pos.x = 0
	m.pos.y = 0
	m.right_flag = false
	m.flag = false
	m.set_Ahead(Dir(3))
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
		fmt.Println("State: player think")
		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		m.input_msg = data[0].([]cell_msg)
		m.executor.Cur_state = "THINK"
	}
}

func (m *think) Output() *system.SysMessage {
	for i := 0; i < len(m.input_msg); i++ {
		direction := m.input_msg[i].dir
		block := m.input_msg[i].block
		if m.ahead.right == direction && m.right_flag == false {
			if block == 0 {
				m.turnRight()
				msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "move")
				output_msg := m.ahead
				msg.Insert(output_msg)
				m.flag = true
				return msg
			} else if block == 1 {
				m.right_flag = true
				m.flag = false
			}
		} else if m.flag == true {
			if m.ahead.front == direction {
				if block == 0 {
					msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "move")
					output_msg := m.ahead
					msg.Insert(output_msg)
					m.flag = true
					return msg
				} else if block == 1 {
					m.turnLeft()
					m.flag = false
				}
			}
		}
	}
	return nil
}

func (m *think) Int_trans() {
	if m.executor.Cur_state == "THINK" && m.flag == true {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "THINK"
	}
}

func (m *think) set_Ahead(ahead Dir) {
	switch ahead {
	case Dir(0):
		m.ahead.front = Dir(0)
		m.ahead.back = Dir(3)
		m.ahead.left = Dir(2)
		m.ahead.right = Dir(1)
	case Dir(3):
		m.ahead.front = Dir(3)
		m.ahead.back = Dir(0)
		m.ahead.left = Dir(1)
		m.ahead.right = Dir(2)
	case Dir(1):
		m.ahead.front = Dir(1)
		m.ahead.back = Dir(2)
		m.ahead.left = Dir(0)
		m.ahead.right = Dir(3)
	case Dir(2):
		m.ahead.front = Dir(2)
		m.ahead.back = Dir(1)
		m.ahead.left = Dir(3)
		m.ahead.right = Dir(0)
	}
}

func (m *think) turnLeft() {
	switch m.ahead.front {
	case Dir(0):
		m.set_Ahead(Dir(2))
	case Dir(3):
		m.set_Ahead(Dir(1))
	case Dir(1):
		m.set_Ahead(Dir(0))
	case Dir(2):
		m.set_Ahead(Dir(3))
	}
}

func (m *think) turnRight() {
	switch m.ahead.front {
	case Dir(0):
		m.set_Ahead(Dir(1))
	case Dir(3):
		m.set_Ahead(Dir(2))
	case Dir(1):
		m.set_Ahead(Dir(3))
	case Dir(2):
		m.set_Ahead(Dir(0))
	}
}

func (m *think) set_position(x int, y int) {
	m.pos.x = x
	m.pos.y = y
}

func (m *think) get_position() pos {
	return m.pos
}
