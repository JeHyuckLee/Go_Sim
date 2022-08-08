package main

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
)

func main() {
	se := executor.NewSysSimulator()
	se.Register_engine("gosim", "VIRTURE_TIME", 1)
	sim := se.Get_engine("gosim")
	sim.Behaviormodel.CoreModel.Insert_input_port("start")

	//player 결합모델
	player := create_player(0, definition.Infinite, "player", "gosim", 1, 1)
	sim.Register_entity(player.am_move.executor)
	sim.Register_entity(player.am_think.executor)
	sim.Coupling_relation(nil, "start", player.am_move.executor, "start")
	sim.Coupling_relation(player.am_think.executor, "move", player.am_move.executor, "think")

	//맵크기
	width := 100
	heigth := 100

	cell := CM_cell(0, definition.Infinite, "cell", "gosim", width, heigth)

	//결합모델 cell 만들기
	sim.Register_entity(cell.am_check.executor)
	sim.Register_entity(cell.am_in.executor)
	sim.Coupling_relation(cell.am_in.executor, "check", cell.am_check.executor, "check")

	//cell player 연결
	sim.Coupling_relation(player.am_move.executor, "in", cell.am_in.executor, "in")
	sim.Coupling_relation(cell.am_check.executor, "player", player.am_think.executor, "player")

	sim.Insert_external_event("start", nil, 0)
	sim.Simulate(definition.Infinite)
}
