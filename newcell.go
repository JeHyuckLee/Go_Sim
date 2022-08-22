package main

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
	"evsim_golang/system"
	"fmt"
	"os"
)

type cm_cell struct {
	am_in    *cell_in
	am_check *cell_check
}

func CM_cell(instance_time, destruct_time float64, name, engine_name string, ix, iy int) *cm_cell {

	cell := cm_cell{}
	cell.am_in = AM_cellIn(instance_time, destruct_time, name, engine_name)
	cell.am_check = AM_cellcheck(instance_time, destruct_time, name, engine_name, ix, iy)

	return &cell
}

type cell_in struct {
	executor *executor.BehaviorModelExecutor
	pos      pos
	msg      *system.SysMessage
}

func AM_cellIn(instance_time, destruct_time float64, name, engine_name string) *cell_in {
	m := &cell_in{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	//state
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("IN", 0)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("in")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("check")

	return m
}

func (m *cell_in) Ext_trans(port string, msg *system.SysMessage) {
	//player가 해당 셀에 왔음
	if port == "in" {
		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		m.pos = data[0].(pos)

		m.executor.Cur_state = "IN"
	}

}

func (m *cell_in) Output() *system.SysMessage {
	//check 에게 출력을 보내서 동작시킴
	fmt.Println("State: cell in")
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "check")
	msg.Insert(m.pos)

	return msg
}

func (m *cell_in) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "IN" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IDLE"
	}
}

type cell_check struct {
	executor *executor.BehaviorModelExecutor
	cell     [][]int
	cell2    [][]int
	pos      pos
	cellmsg  []cell_msg
}

func AM_cellcheck(instance_time, destruct_time float64, name, engine_name string, x, y int) *cell_check {
	m := &cell_check{}
	m.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
	m.executor.AbstractModel = m

	m.cellmsg = make([]cell_msg, 4)
	m.cell = create_map(x, y)
	m.cell2 = create_map(x, y)
	m.cellmsg[0].dir = Dir_Nort
	m.cellmsg[1].dir = Dir_East
	m.cellmsg[2].dir = Dir_West
	m.cellmsg[3].dir = Dir_South
	//state
	m.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	m.executor.Behaviormodel.Insert_state("CHECK", 0)
	m.executor.Init_state("IDLE")

	//port
	m.executor.Behaviormodel.CoreModel.Insert_input_port("check")
	m.executor.Behaviormodel.CoreModel.Insert_output_port("player")

	return m
}

func (m *cell_check) Ext_trans(port string, msg *system.SysMessage) {
	//player가 해당 셀에 왔음
	if port == "check" {
		file1, _ := os.Create("map.csv")

		m.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		m.pos = data[0].(pos)
		m.cell2[m.pos.x][m.pos.y] = +2
		for i := 0; i < 100; i++ {
			for j := 0; j < 100; j++ {
				fmt.Fprint(file1, m.cell2[j][i], ",")
			}
			fmt.Fprintln(file1)
		}
		file1.Close()
		fmt.Println("State: cell check")
		m.executor.Cur_state = "CHECK"
	}

}

func (m *cell_check) Output() *system.SysMessage {
	//check 에게 출력을 보내서 동작시킴
	//시작지점을 1.1 로 하고 주변을 벽으로 해야할듯
	m.cellmsg[0].pos.x = m.pos.x
	m.cellmsg[1].pos.x = m.pos.x + 1
	m.cellmsg[2].pos.x = m.pos.x - 1
	m.cellmsg[3].pos.x = m.pos.x

	m.cellmsg[0].pos.y = m.pos.y - 1
	m.cellmsg[1].pos.y = m.pos.y
	m.cellmsg[2].pos.y = m.pos.y
	m.cellmsg[3].pos.y = m.pos.y + 1

	// 1 = block
	for i := 0; i < 4; i++ {
		m.cellmsg[i].block = m.cell[m.cellmsg[i].pos.y][m.cellmsg[i].pos.x]
	}
	msg := system.NewSysMessage(m.executor.Behaviormodel.CoreModel.Get_name(), "player")
	msg.Insert(m.cellmsg)

	return msg
}

func (m *cell_check) Int_trans() {
	//상태변화
	if m.executor.Cur_state == "CHECK" {
		m.executor.Cur_state = "IDLE"
	} else {
		m.executor.Cur_state = "IDLE"
	}
}
