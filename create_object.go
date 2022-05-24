// package main

// import (
// 	"evsim_golang/definition"
// 	"evsim_golang/executor"
// 	"evsim_golang/system"
// 	"fmt"
// 	"os"
// 	"time"
// )

// type Object struct {
// 	executor *executor.BehaviorModelExecutor
// }

// func (ob *Object) Ext_trans(port string, msg *system.SysMessage) {
// 	// fmt.Println("[object][Ext_trans]")
// 	if port == "object" {
// 		// fmt.Println("[object][In]:")
// 		ob.executor.Cur_state = "MOVE"
// 	}
// }

// func (ob *Object) Output() *system.SysMessage {
// 	// fmt.Println("[object][Output]")
// 	startTime := time.Now()
// 	const NUM_OBJECT int = 100

// 	for i := 1; i <= NUM_OBJECT; i++ {
// 		file, err := os.Create("C:\\evsim_golang\\src\\object\\" + fmt.Sprintf("object%d", i) + ".txt")
// 		if err != nil {
// 			panic(err)
// 		}
// 		defer file.Close()

// 		for j := 1; j <= NUM_OBJECT; j++ {
// 			_, err = file.WriteString(fmt.Sprintf("%d\n", j))
// 			if err != nil {
// 				panic(err)
// 			}
// 		}
// 	}
// 	Total_time := time.Since(startTime)
// 	fmt.Printf("Total Time: %s\n", Total_time)
// 	msg := system.NewSysMessage(ob.executor.Behaviormodel.CoreModel.Get_name(), "obj")
// 	return msg
// }

// func (ob *Object) Int_trans() {
// 	// fmt.Println("[object][Int_trans]")
// 	if ob.executor.Cur_state == "MOVE" {
// 		ob.executor.Cur_state = "IDLE"
// 	} else {
// 		ob.executor.Cur_state = "MOVE"
// 	}
// }

// func NewObject() *Object {
// 	obj := Object{}
// 	obj.executor = executor.NewExecutor(0, definition.Infinite, "obj", "sname")
// 	obj.executor.AbstractModel = &obj
// 	obj.executor.Init_state("IDLE")
// 	obj.executor.Behaviormodel.Insert_state("IDLE", definition.Infinite)
// 	obj.executor.Behaviormodel.Insert_state("MOVE", 1)
// 	obj.executor.Behaviormodel.CoreModel.Insert_input_port("object")
// 	return &obj
// }

// func main() {
// 	se := executor.NewSysSimulator()
// 	se.Register_engine("sname", "REAL_TIME", 1)
// 	sim := se.Get_engine("sname")
// 	sim.Behaviormodel.CoreModel.Insert_input_port("object")

// 	obj := NewObject()

// 	sim.Register_entity(obj.executor)

// 	sim.Coupling_relation(nil, "object", obj.executor, "object")
// 	sim.Insert_external_event("object", nil, 0)
// 	sim.Simulate(definition.Infinite)

// }
// 