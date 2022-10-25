package main

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

func NewGenerator(instance_time, destruct_time float64, name, engine_name string) *Generator {
	gen := Generator{}
	gen.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
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

func (p *Processor) Output() *system.SysMessage {
	msg := system.NewSysMessage(p.executor.Behaviormodel.CoreModel.Get_name(), "generator")
	msg.Insert("ack")
	return msg
}

func NewProcessor(instance_time, destruct_time float64, name, engine_name string) *Processor {
	pro := &Processor{}
	pro.executor = executor.NewExecutor(instance_time, destruct_time, name, engine_name)
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
		gen := NewGenerator(0, definition.Infinite, "Gen", "sname")
		pro := NewProcessor(0, definition.Infinite, "Proc", "sname")
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
