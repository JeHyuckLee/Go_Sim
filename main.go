package main

import (
	"evsim_golang/definition"
	"evsim_golang/executor"
	"evsim_golang/system"
	"fmt"
	"runtime"
	"time"
)

type Generator struct {
	executor *executor.BehaviorModelExecutor
	msg_list []interface{}
}

func (g *Generator) Ext_trans(port string, msg *system.SysMessage) {

	//fmt.Println("ext_trans")
	if port == "start" {
		fmt.Println("[gen][in]:", time.Now())

		g.executor.Cur_state = "MOVE"
	}
	fmt.Println("\n gen_ext_trans :", time.Since(executor.Start_time))
}

func (g *Generator) Int_trans() {
	//fmt.Println("int_trans")
	if g.executor.Cur_state == "MOVE" && g.msg_list == nil {
		g.executor.Cur_state = "IDLE"
	} else {
		g.executor.Cur_state = "MOVE"
	}
	fmt.Println("\n get_int_trans :", time.Since(executor.Start_time))
}

func (g *Generator) Output() *system.SysMessage {
	//fmt.Println("output")
	msg := system.NewSysMessage(g.executor.Behaviormodel.CoreModel.Get_name(), "process")
	fmt.Println("[gen][out]:", time.Now())
	msg.Insert(g.msg_list[0])
	g.msg_list = remove(g.msg_list, 0)
	fmt.Println("\n gen_outPut :", time.Since(executor.Start_time))
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
	for i := 0; i < 10; i++ {
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
		fmt.Println("[proc][in]", time.Now())
		p.executor.Cancel_rescheduling()
		data := msg.Retrieve()
		p.msg_list = append(p.msg_list, data...)
		p.executor.Cur_state = "PROCESS"
	}
	fmt.Println("\n pro_exttrans :", time.Since(executor.Start_time))
}

func (p *Processor) Int_trans() {
	//fmt.Println("int_trans")
	if p.executor.Cur_state == "PROCESS" {
		p.executor.Cur_state = "IDLE"
	} else {
		p.executor.Cur_state = "IDLE"
	}
	fmt.Println("\n pro_int_trans :", time.Since(executor.Start_time))
}

func (p Processor) Output() *system.SysMessage {
	//fmt.Println("output")
	fmt.Println("[proc][out]", time.Now())
	fmt.Println(p.msg_list...)

	fmt.Println("\nproc_output :", time.Since(executor.Start_time))
	return nil
}

func NewProcessor() *Processor {
	pro := &Processor{}
	pro.executor = executor.NewExecutor(0, definition.Infinite, "Proc", "sname")
	pro.executor.AbstractModel = pro
	pro.executor.Init_state("IDLE")
	pro.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
	pro.executor.Behaviormodel.Insert_state("PROCESS", 2)
	pro.executor.Behaviormodel.CoreModel.Insert_input_port("PROCESS")

	return pro
}

func main() {
	fmt.Println("start", time.Now())
	executor.Start_time = time.Now()
	runtime.GOMAXPROCS(8)
	se := executor.NewSysSimulator()
	se.Register_engine("sname", "REAL_TIME", 1)
	sim := se.Get_engine("sname")
	sim.Behaviormodel.CoreModel.Insert_input_port("start")
	for i := 0; i < 1; i++ {
		gen := NewGenerator()
		pro := NewProcessor()
		sim.Register_entity(gen.executor)
		sim.Register_entity(pro.executor)
		sim.Coupling_relation(nil, "start", gen.executor, "start")
		sim.Coupling_relation(gen.executor, "process", pro.executor, "process")
	}

	sim.Insert_external_event("start", nil, 0)
	sim.Simulate(definition.Infinite)

}

func remove(slice []interface{}, s int) []interface{} {
	return append(slice[:s], slice[s+1:]...)
}
