package main

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
	"fmt"
)

//player 는 처음 0.0 에게 입력을 보냄 입력을 받은 0.0은 연결된 0.1, 1.0 에게
//입력을 보내고 블럭인지 아닌지에 대한 정보를 받아와서 플레이어 에게 알려줌
//입력을 받은 player 는 어떻게 움직일지 판단한 후에 다음 블럭에게 입력을 보냄

func main() {
	se := executor.NewSysSimulator()
	se.Register_engine("gosim", "VIRTURE_TIME", 1)
	sim := se.Get_engine("gosim")
	sim.Behaviormodel.CoreModel.Insert_input_port("start")

	//player 결합모델
	am_move := AM_move(0, definition.Infinite, "move", "gosim")
	am_think := AM_think(0, definition.Infinite, "think", "gosim")
	sim.Register_entity(am_move.executor)
	sim.Register_entity(am_think.executor)
	sim.Coupling_relation(nil, "start", am_move.executor, "start")
	sim.Coupling_relation(am_think.executor, "move", am_move.executor, "think")

	//맵크기
	width := 10
	heigth := 10
	// cell 끼리 연결 을 위해 만든 슬라이스
	cell_Check := make([][]*check, heigth)
	cell_Info := make([][]*cell_info, heigth)
	for i := 0; i < heigth; i++ {
		cell_Check[i] = make([]*check, width)
		cell_Info[i] = make([]*cell_info, width)

		for j := 0; j < width; j++ {
			//cell의 원자모델 들 생성
			am_check := AM_check(0, definition.Infinite, "check", "gosim")
			am_in := AM_cellIn(0, definition.Infinite, "in", "gosim")
			am_out := AM_cellOut(0, definition.Infinite, "out", "gosim")
			am_info := AM_cellInfo(0, definition.Infinite, "info", "gosim", j, i)
			n := fmt.Sprintf("{%d,%d}", j, i)
			am_check.executor.Behaviormodel.CoreModel.Set_name(n)
			am_out.executor.Behaviormodel.CoreModel.Set_name(n)
			am_in.executor.Behaviormodel.CoreModel.Set_name(n)
			am_info.executor.Behaviormodel.CoreModel.Set_name(n)
			am_info.set_position(j, i)
			//결합모델 cell 만들기
			sim.Register_entity(am_check.executor)
			sim.Register_entity(am_out.executor)
			sim.Register_entity(am_in.executor)
			sim.Register_entity(am_info.executor)

			sim.Coupling_relation(am_in.executor, "check", am_check.executor, "in")
			// sim.Coupling_relation(am_check.executor, "out", am_out.executor, "check")

			//player 와 cell 의 연결
			am_move.insert_Player_Output_Port(n)
			am_in.insert_cell_Input_Port(n)

			sim.Coupling_relation(am_move.executor, n, am_in.executor, n)
			sim.Coupling_relation(am_out.executor, "player", am_think.executor, "think")
			cell_Check[i][j] = am_check
			cell_Info[i][j] = am_info

		}

	}

	//cell과 cell의 연결
	for i := 0; i < heigth; i++ {
		for j := 0; j < width; j++ {

			if i != 0 {
				sim.Coupling_relation(cell_Check[i][j].executor, "south", cell_Info[i-1][j].executor, "north")
				sim.Coupling_relation(cell_Info[i-1][j].executor, "north", cell_Check[i][j].executor, "south")
				cell_Check[i][j].con_list[3] = true // n =0, e =1, w =2, s= 3
			}
			if i != heigth-1 {
				sim.Coupling_relation(cell_Check[i][j].executor, "north", cell_Info[i+1][j].executor, "south")
				sim.Coupling_relation(cell_Info[i+1][j].executor, "south", cell_Check[i][j].executor, "north")
				cell_Check[i][j].con_list[0] = true
			}

			if j != 0 {
				sim.Coupling_relation(cell_Check[i][j].executor, "west", cell_Info[i][j-1].executor, "east")
				sim.Coupling_relation(cell_Info[i][j-1].executor, "east", cell_Check[i][j].executor, "west")
				cell_Check[i][j].con_list[2] = true
			}
			if j != width-1 {
				sim.Coupling_relation(cell_Check[i][j].executor, "east", cell_Info[i][j+1].executor, "west")
				sim.Coupling_relation(cell_Info[i][j+1].executor, "west", cell_Check[i][j].executor, "east")
				cell_Check[i][j].con_list[1] = true
			}
		}
	}

	//나중에 맵을 2차원 배열로 만들면 이렇게하면될듯
	// for i := 0; i < heigth; i++ {
	// 	for j := 0; j < width; j++ {
	// 		cell_Info[i][j].set_block() = map_array[i][j]
	// 	}
	// }

	//시작
	sim.Insert_external_event("start", nil, 0)
	sim.Simulate(definition.Infinite)
}
