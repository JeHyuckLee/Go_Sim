package main

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
	"evsim_golang/system"
	"fmt"
)

//cell의 원자모델
type cellOut struct {
	executor  *executor.BehaviorModelExecutor
	msg_list  []interface{}
	cell_list map[Dir]pos
	cell_msg
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

func (m *cellOut) Ext_trans(port string, msg *system.SysMessage) {
	//check 에게 정보를 받음
	if port == "check" {
		m.executor.Cancel_rescheduling()
		// data := msg.Retrieve()

	}
}

func (m *cellOut) Output() *system.SysMessage {
	//player 에게 전송
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "player")

	return msg
}

func (m *cellOut) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "OUT" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "OUT"
	}
}

//cell의 원자모델
type cellIn struct {
	executor    *executor.BehaviorModelExecutor
	player_list []interface{}
}

func AM_cellIn(instance_time, destruct_time float64, name, engine_name string) *cellIn {
	m := cellIn{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = &m

	//state
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("IN", 0)
	m.executor.Init_state("IDLE")

	//port

	m.executor.Behaviormodel.CoreModel.Insert_output_port("check")

	return &m
}

func (m *cellIn) insert_cell_Input_Port(port string) {
	m.executor.Behaviormodel.CoreModel.Insert_input_port(port)
}

func (m *cellIn) Ext_trans(port string, msg *system.SysMessage) {
	//player가 해당 셀에 왔음
	if port == m.executor.Behaviormodel.CoreModel.Get_name() {
		fmt.Println("Player IN: ", m.executor.Behaviormodel.CoreModel.Get_name())
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
	// m.player_list = remove(m.player_list, 0)

	return msg
}

func (m *cellIn) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "IN" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IN"
	}
}

//cell의 원자모델
type check struct {
	executor   *executor.BehaviorModelExecutor
	block_list []interface{}
	block      bool
	con_count  int
	count      int
	checking   bool
	Nflag      bool
	Sflag      bool
	Eflag      bool
	Wflag      bool
	out        bool
	con_list   map[Dir]bool //n =0 e = 1 w = 2 s = 3
	output     []cell_msg
	out_dir    Dir
	out_port   string
}

func AM_check(instance_time, destruct_time float64, name, engine_name string) *check {
	//맵 모델
	m := check{}
	m.checking = false
	m.con_list = make(map[Dir]bool)
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = &m
	m.block = false
	m.Nflag = false
	m.Sflag = false
	m.Eflag = false
	m.Wflag = false

	//state

	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	//OUT으로부터 입력이 오면 IDLE -> CHECK
	m.executor.Behaviormodel.Insert_state("CHECK", 0)
	//Chcek한 셀 정보 넘기기
	m.executor.Behaviormodel.Insert_state("SET", 0)
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

func (m *check) flag_initialize() {
	//con_list 가 true 이면 연결되어있음
	for k, v := range m.con_list {
		if v == true {
			m.count++
			m.con_count++
		}
		if k == 0 {
			m.Nflag = v
		} else if k == 1 {
			m.Eflag = v
		} else if k == 2 {
			m.Wflag = v
		} else if k == 3 {
			m.Sflag = v
		}
	}
}
func (m *check) Ext_trans(port string, msg *system.SysMessage) {
	//in에게 입력을 받으면 check 상태로 가고 연결된 애들에게 입력을 보냄
	//NEWS 포트로 입력을 받으면 SET 상태로 가고 각 셀 정보 저장
	//전부 탐색시 OUT 상태로 가고 out에 정보 전달

	//player 가 셀에 들어왔다
	if port == "in" {
		m.flag_initialize()
		m.executor.Cur_state = "CHECK"
	} else { //입력보냈던 셀 로 부터 정보를 다시 받음
		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		msg := data[0].(cell_msg)
		m.output = append(m.output, msg)
		m.count++
		if m.count == m.con_count {
			m.executor.Cur_state = "OUT"
		}
		m.executor.Cur_state = "SET"
	}
}

func (m *check) Output() *system.SysMessage {
	//in에게 입력을 받으면 NEWS 포트중 연결된 포트로 출력
	//NEWS포트 로 입력이 들어오면 입력된 정보를 상대방 셀 에게 전송
	if m.executor.Cur_state == "CHECK" {
		//true 이면 연결되어있음
		if m.Nflag == true {
			msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "north")
			m.Nflag = false
			m.count--
			return msg
		} else if m.Sflag == true {
			msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "south")
			m.Sflag = false
			m.count--
			return msg
		} else if m.Eflag == true {
			msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "east")
			m.Eflag = false
			m.count--
			return msg
		} else if m.Wflag == true {
			msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "west")
			m.Wflag = false
			m.count--
			return msg
		}
	}

	if m.executor.Cur_state == "OUT" {
		msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "out")
		msg.Insert(m.output)
		return msg
	}

	return nil
}

func (m *check) Int_trans() {
	//상태천이
	if m.executor.Cur_state == "CHECK" && m.count == 0 {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IDLE"
	}
}

type cell_info struct {
	executor *executor.BehaviorModelExecutor
	pos
	cell_msg
	out_dir  Dir
	out_port string
	block    bool
}

func AM_cellInfo(instance_time, destruct_time float64, name, engine_name string, ix, iy int) *cell_info {
	m := cell_info{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.set_position(ix, iy)
	m.executor.AbstractModel = &m

	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("INFO", 1)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("north")
	m.executor.Behaviormodel.CoreModel.Insert_input_port("south")
	m.executor.Behaviormodel.CoreModel.Insert_input_port("east")
	m.executor.Behaviormodel.CoreModel.Insert_input_port("west")

	m.executor.Behaviormodel.CoreModel.Insert_output_port("north")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("south")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("east")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("west")

	return &m

}

func (m *cell_info) set_position(x int, y int) {
	m.pos.x = x
	m.pos.y = y
}

func (m *cell_info) get_position() pos {
	return m.pos
}

func (m *cell_info) set_block(b bool) {
	m.block = b
}
func (m *cell_info) get_block() bool {
	return m.block
}

func (m *cell_info) Ext_trans(port string, msg *system.SysMessage) {

	if port == "north" {

		m.executor.Cur_state = "out"
		m.out_dir = 3     //입력이 들어온 셀을 기준으로 했을때 자신의 위치
		m.out_port = port // 입력이 들어온 셀로 다시 보내줌

	} else if port == "east" {

		m.executor.Cur_state = "out"
		m.out_dir = 2
		m.out_port = port

	} else if port == "west" {

		m.executor.Cur_state = "out"
		m.out_dir = 1
		m.out_port = port

	} else if port == "south" {

		m.executor.Cur_state = "out"
		m.out_dir = 0
		m.out_port = port
	}
}

func (m *cell_info) Output() *system.SysMessage {
	//check 에게 출력을 보내서 동작시킴
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), m.out_port)
	m.cell_msg.dir = m.out_dir
	m.cell_msg.pos = m.get_position()
	m.cell_msg.block = m.get_block()
	msg.Insert(m.cell_msg)

	return msg
}

func (m *cell_info) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "INFO" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "INFO"
	}
}
