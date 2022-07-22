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
	ac_move := AC_move()
	ac_think := AC_think()
	sim.Register_entity(ac_move.executor)
	sim.Register_entity(ac_think.executor)
	sim.Coupling_relation(nil, "start", ac_move.executor, "start")
	sim.Coupling_relation(ac_think.executor, "move", ac_move.executor, "think")

	//맵크기
	width := 100
	heigth := 100
	// cell 끼리 연결 을 위해 만든 슬라이스
	cell := make([][]*executor.BehaviorModelExecutor, heigth)

	for i := 0; i < heigth; i++ {
		cell[i] = make([]*executor.BehaviorModelExecutor, width)

		for j := 0; j < width; j++ {
			//cell의 원자모델 들 생성
			ac_check := AC_check()
			ac_in := AC_cellin()
			ac_out := AC_cellout()
			n := fmt.Sprintf("{%n,%n}", j, i)
			ac_check.executor.Behaviormodel.CoreModel.Set_name(n)
			ac_out.executor.Behaviormodel.CoreModel.Set_name(n)
			ac_in.executor.Behaviormodel.CoreModel.Set_name(n)

			//결합모델 cell 만들기
			sim.Register_entity(ac_check.executor)
			sim.Register_entity(ac_out.executor)
			sim.Register_entity(ac_in.executor)
			sim.Coupling_relation(ac_in.executor, "check", ac_check.executor, "in")
			sim.Coupling_relation(ac_check.executor, "out", ac_out.executor, "check")

			//player 와 cell 의 연결
			sim.Coupling_relation(ac_move.executor, "cell", ac_in.executor, "cell")
			sim.Coupling_relation(ac_out.executor, "player", ac_think.executor, "think")
			cell[i][j] = ac_check.executor
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
