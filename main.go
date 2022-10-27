package main

// import (
// 	"evsim_golang/executor"
// 	"fmt"
// 	"runtime"
// 	"time"
// )

// func main() {

// 	fmt.Println("start", time.Now())
// 	executor.Start_time = time.Now()
// 	runtime.GOMAXPROCS(8)
// 	se := executor.NewSysSimulator()
// 	se.Register_engine("sname", "VIRTURE_TIME", 1)
// 	sim := se.Get_engine("sname")
// 	sim.Behaviormodel.CoreModel.Insert_input_port("start")

// 	coop := CM_coop(0, 150, "coop", "sname", 100)
// 	sim.Register_entity(coop.am_management.executor)
// 	sim.Register_entity(coop.am_shipment.executor)
// 	sim.Register_entity(coop.am_ware.executor)

// 	buyer := AM_buyer(0, 150, "buyer", "sname", 120)
// 	sim.Register_entity(buyer.executor)

// 	sim.Coupling_relation(buyer.executor, "buy", coop.am_shipment.executor, "shipment")

// 	for i := 0; i < 5; i++ {
// 		member := CM_coopMember(0, 150, "member", "sname", 50, 100, 100)
// 		sim.Register_entity(member.am_harvest.executor)
// 		sim.Register_entity(member.am_seed.executor)
// 		sim.Register_entity(member.am_ship.executor)
// 		sim.Coupling_relation(member.am_seed.executor, "harvest", member.am_harvest.executor, "harvest")
// 		sim.Coupling_relation(member.am_harvest.executor, "shipment", member.am_ship.executor, "shipment")
// 		sim.Coupling_relation(member.am_ship.executor, "in", coop.am_ware.executor, "in")
// 		sim.Coupling_relation(nil, "start", member.am_seed.executor, "seeding")
// 		// sim.Register_parallel_entity(gen.executor)
// 	}
// 	sim.Coupling_relation(nil, "start", buyer.executor, "start")

// 	sim.Insert_external_event("start", nil, 0)
// 	sim.Simulate(100)

// }
import (
	"log"
	"time"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	zctx, _ := zmq.NewContext()

	s, _ := zctx.NewSocket(zmq.REP)
	s.Bind("tcp://*:5555")

	for {
		// Wait for next request from client
		msg, _ := s.Recv(0)
		log.Printf("Received %s\n", msg)

		// Do some 'work'
		time.Sleep(time.Second * 1)

		// Send reply back to client
		s.Send("World", 0)
	}
}
