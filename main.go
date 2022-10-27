package main

import (
	"evsim_golang/executor"
	"fmt"
	"runtime"
	"time"
)

type member_infor struct {
	aver      float64
	Member_id int
	std       float64
	cnt       int
	area      int
}

type buyer_infor struct {
	aver float64
	cnt  int
	std  float64
}

func main() {

	// 데이터 불러오기
	db := GetConnector()
	defer db.Close()

	err := db.Ping()
	if err != nil {
		panic(err)
	}
	var memberList []member_infor
	results, err := db.Query("SELECT i.Member_id ,AVG(i.Warehousing_amount) as aver,STD(i.Warehousing_amount) as std ,COUNT(i.Warehousing_amount) as cnt,cm.Member_area FROM inventory i left join coopMember cm on i.Member_id = cm.Member_id GROUP BY i.Member_id ")
	for results.Next() {
		var member_infor member_infor
		err = results.Scan(&member_infor.Member_id, &member_infor.aver, &member_infor.std, &member_infor.cnt, &member_infor.area)
		if err != nil {
			panic(err.Error())
		}
		memberList = append(memberList, member_infor)
	}

	var buyer_inf buyer_infor
	results, err = db.Query("SELECT AVG(s.Shipment_amount) as aver, COUNT(s.Shipment_amount) as cnt, STD(s.Shipment_amount) as std FROM Shipment s")
	for results.Next() {
		var buyer buyer_infor
		err = results.Scan(&buyer.aver, &buyer.cnt, &buyer.std)
		if err != nil {
			panic(err.Error())
		}
		buyer_inf = buyer
	}

	// 엔진설정
	fmt.Println("start", time.Now())
	executor.Start_time = time.Now()
	runtime.GOMAXPROCS(8)
	se := executor.NewSysSimulator()
	se.Register_engine("sname", "VIRTURE_TIME", 1)
	sim := se.Get_engine("sname")
	sim.Behaviormodel.CoreModel.Insert_input_port("start")

	coop := CM_coop(0, 150, "coop", "sname", 100)
	sim.Register_entity(coop.am_management.executor)
	sim.Register_entity(coop.am_shipment.executor)
	sim.Register_entity(coop.am_ware.executor)

	buy := rand_crop(buyer_inf.aver, buyer_inf.std)
	buyer := AM_buyer(0, 150, "buyer", "sname", int(buy)*buyer_inf.cnt/2)
	sim.Register_entity(buyer.executor)

	sim.Coupling_relation(buyer.executor, "buy", coop.am_shipment.executor, "shipment")

	for _, v := range memberList {
		crop := rand_crop(v.aver, v.std)
		member := CM_coopMember(0, 150, "member", "sname", v.area, crop, 50)
		sim.Register_entity(member.am_harvest.executor)
		sim.Register_entity(member.am_seed.executor)
		sim.Register_entity(member.am_ship.executor)
		sim.Coupling_relation(member.am_seed.executor, "harvest", member.am_harvest.executor, "harvest")
		sim.Coupling_relation(member.am_harvest.executor, "shipment", member.am_ship.executor, "shipment")
		sim.Coupling_relation(member.am_ship.executor, "in", coop.am_ware.executor, "in")
		sim.Coupling_relation(nil, "start", member.am_seed.executor, "seeding")
		// sim.Register_parallel_entity(gen.executor)
	}
	sim.Coupling_relation(nil, "start", buyer.executor, "start")

	sim.Insert_external_event("start", nil, 0)
	sim.Simulate(365)
}
