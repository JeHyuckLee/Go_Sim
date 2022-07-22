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
	am_move := AM_move()
	am_think := AM_think()
	sim.Register_entity(am_move.executor)
	sim.Register_entity(am_think.executor)
	sim.Coupling_relation(nil, "start", am_move.executor, "start")
	sim.Coupling_relation(am_think.executor, "move", am_move.executor, "think")

	//맵크기
	width := 100
	heigth := 100
	// cell 끼리 연결 을 위해 만든 슬라이스
	cell := make([][]*executor.BehaviorModelExecutor, heigth)

	for i := 0; i < heigth; i++ {
		cell[i] = make([]*executor.BehaviorModelExecutor, width)

		for j := 0; j < width; j++ {
			//cell의 원자모델 들 생성
			am_check := AM_check()
			am_in := AM_cellin()
			am_out := AM_cellout()
			n := fmt.Sprintf("{%n,%n}", j, i)
			am_check.executor.Behaviormodel.CoreModel.Set_name(n)
			am_out.executor.Behaviormodel.CoreModel.Set_name(n)
			am_in.executor.Behaviormodel.CoreModel.Set_name(n)

			//결합모델 cell 만들기
			sim.Register_entity(am_check.executor)
			sim.Register_entity(am_out.executor)
			sim.Register_entity(am_in.executor)
			sim.Coupling_relation(am_in.executor, "check", am_check.executor, "in")
			sim.Coupling_relation(am_check.executor, "out", am_out.executor, "check")

			//player 와 cell 의 연결
			sim.Coupling_relation(am_move.executor, "cell", am_in.executor, "cell")
			sim.Coupling_relation(am_out.executor, "player", am_think.executor, "think")
			cell[i][j] = am_check.executor
		}

	}

	//cell과 cell의 연결
	for i := 0; i < heigth; i++ {
		for j := 0; j < width; j++ {

			if i != 0 {
				sim.Coupling_relation(cell[i][j], "south", cell[i-1][j], "north")
			}
			if i != heigth-1 {
				sim.Coupling_relation(cell[i][j], "north", cell[i+1][j], "south")
			}

			if j != 0 {
				sim.Coupling_relation(cell[i][j], "west", cell[i][j-1], "east")
			}
			if j != width-1 {
				sim.Coupling_relation(cell[i][j], "east", cell[i][j+1], "west")
			}
		}
	}

	//시작
	sim.Insert_external_event("start", nil, 0)
	sim.Simulate(definition.Infinite)
}
