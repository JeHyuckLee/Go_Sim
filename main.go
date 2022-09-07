package main

// func main() {
// 	se := executor.NewSysSimulator()
// 	se.Register_engine("gosim", "VIRTURE_TIME", 1)
// 	sim := se.Get_engine("gosim")
// 	sim.Behaviormodel.CoreModel.Insert_input_port("start")
// 	width := 100
// 	heigth := 100

// 	for i := 0; i < 5; i++ {
// 		cell := CM_cell(0, definition.Infinite, "cell", "gosim", width, heigth)

// 		sim.Register_entity(cell.am_check.executor)
// 		sim.Register_entity(cell.am_in.executor)
// 		sim.Coupling_relation(cell.am_in.executor, "check", cell.am_check.executor, "check")
// 		player := create_player(0, definition.Infinite, "player", "gosim", 1, 1)
// 		sim.Register_entity(player.am_move.executor)
// 		sim.Register_entity(player.am_think.executor)
// 		sim.Coupling_relation(nil, "start", player.am_move.executor, "start")
// 		sim.Coupling_relation(player.am_think.executor, "move", player.am_move.executor, "think")
// 		sim.Coupling_relation(player.am_move.executor, "in", cell.am_in.executor, "in")
// 		sim.Coupling_relation(cell.am_check.executor, "player", player.am_think.executor, "player")

// 	}

// 	sim.Insert_external_event("start", nil, 0)
// 	sim.Simulate(definition.Infinite)
// }

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
	"evsim_golang/system"
	"fmt"
	"time"
)

type Generator struct {
	executor *executor.BehaviorModelExecutor
	msg_list []interface{}
}

func (g *Generator) Ext_trans(port string, msg *system.SysMessage) {

	//fmt.Println("ext_trans")
	if port == "start" {
		if len(g.msg_list) != 0 {
			g.executor.Cur_state = "MOVE"
		}

	}
}

func (g *Generator) Int_trans() {
	//fmt.Println("int_trans")
	if g.executor.Cur_state == "MOVE" && len(g.msg_list) == 0 {
		g.executor.Cur_state = "IDLE"
	} else {
		g.executor.Cur_state = "IDLE"
	}
}

func (g *Generator) Output() *system.SysMessage {
	//fmt.Println("output")
	msg := system.NewSysMessage(g.executor.Behaviormodel.CoreModel.Get_name(), "process")
	msg.Insert(g.msg_list[0])
	g.msg_list = remove(g.msg_list, 0)
	return msg
}

func NewGenerator() *Generator {
	gen := Generator{}
	gen.executor = executor.NewExecutor(0, definition.Infinite, "Gen", "sname")
	gen.executor.AbstractModel = &gen
	gen.executor.Init_state("IDLE")
	gen.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	gen.executor.Behaviormodel.Insert_state("MOVE", 1)
	gen.executor.Behaviormodel.CoreModel.Insert_input_port("start")
	gen.executor.Behaviormodel.CoreModel.Insert_output_port("process")

	for i := 0; i < 5; i++ {
		gen.msg_list = append(gen.msg_list, i)
	}
	return &gen
}

type Processor struct {
	executor *executor.BehaviorModelExecutor
	msg_list []interface{}
}

func (p *Processor) Ext_trans(port string, msg *system.SysMessage) {
	//fmt.Println("ext_trans")
	if port == "process" {
		p.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		p.msg_list = append(p.msg_list, data...)
		p.executor.Cur_state = "PROCESS"
	}
}

func (p *Processor) Int_trans() {
	//fmt.Println("int_trans")
	if p.executor.Cur_state == "PROCESS" {
		p.executor.Cur_state = "IDLE"
	} else {
		p.executor.Cur_state = "IDLE"
	}
}

func (p Processor) Output() *system.SysMessage {
	msg := system.NewSysMessage(p.executor.Behaviormodel.CoreModel.Get_name(), "generator")
	msg.Insert("ack")
	return msg
}

func NewProcessor() *Processor {
	pro := &Processor{}
	pro.executor = executor.NewExecutor(0, definition.Infinite, "Proc", "sname")
	pro.executor.AbstractModel = pro
	pro.executor.Init_state("IDLE")
	pro.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	pro.executor.Behaviormodel.Insert_state("PROCESS", 1)
	pro.executor.Behaviormodel.CoreModel.Insert_input_port("process")
	pro.executor.Behaviormodel.CoreModel.Insert_output_port("generator")
	return pro
}

func main() {
	fmt.Println("start", time.Now())
	executor.Start_time = time.Now()
	se := executor.NewSysSimulator()
	se.Register_engine("sname", "VIRTURE_TIME", 1)
	sim := se.Get_engine("sname")
	sim.Behaviormodel.CoreModel.Insert_input_port("start")
	for i := 0; i < 5; i++ {
		gen := NewGenerator()
		pro := NewProcessor()
		sim.Register_entity(gen.executor)
		sim.Register_entity(pro.executor)
		sim.Coupling_relation(nil, "start", gen.executor, "start")
		sim.Coupling_relation(gen.executor, "process", pro.executor, "process")
		sim.Coupling_relation(pro.executor, "generator", gen.executor, "start")

		sim.Register_parallel_entity(gen.executor)
		sim.Register_parallel_entity(pro.executor)
	}

	sim.Insert_external_event("start", nil, 0)
	fmt.Println("model_running_time : ", time.Since(executor.Start_time))
	sim.Simulate(definition.Infinite)

}

func remove(slice []interface{}, s int) []interface{} {
	return append(slice[:s], slice[s+1:]...)
}
