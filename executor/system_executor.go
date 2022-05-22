package executor

import (
	"container/heap"
	"errors"
	"evsim_golang/definition"
	"evsim_golang/model"
	"evsim_golang/system"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/gammazero/deque"
	"gopkg.in/getlantern/deepcopy.v1"
)

type SysExecutor struct {
	// sysObject     *system.SysObject
	Behaviormodel      *model.Behaviormodel
	dmc                *DefaultMessageCatcher
	global_time        float64
	target_time        float64
	time_step          float64
	EXTERNAL_SRC       string
	EXTERNAL_DST       string
	simulation_mode    int
	min_schedule_item  []*BehaviorModelExecutor
	input_event_queue  input_heap
	output_event_queue []*o_event_queue
	sim_mode           string
	waiting_obj_map    map[float64][]*BehaviorModelExecutor
	active_obj_map     map[float64]*BehaviorModelExecutor
	port_map           map[Object][]Object
	sim_init_time      time.Time
}

// deque sort

type schedule_item []*BehaviorModelExecutor

func (s schedule_item) Len() int {
	return len(s)
}
func (s schedule_item) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s schedule_item) Less(i, j int) bool {
	return s[i].Get_req_time() < s[j].Get_req_time()
}

type Object struct {
	object *BehaviorModelExecutor
	port   string
}

type i_event_queue struct {
	time float64
	msg  *system.SysMessage
}

type o_event_queue struct {
	time     float64
	msg_list interface{}
}

type input_heap []i_event_queue

func (eq input_heap) Len() int {
	return len(eq)
}

func (eq input_heap) Less(i, j int) bool {
	return false
}

func (eq input_heap) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
}

func (eq *input_heap) Push(elem interface{}) {
	*eq = append(*eq, elem.(i_event_queue))
}

func (eq *input_heap) Pop() interface{} {
	old := *eq
	n := len(old)
	elem := old[n-1]
	*eq = old[0 : n-1]

	return elem
}

//생성자
func NewSysExecutor(_time_step float64, _sim_name, _sim_mode string) *SysExecutor {
	se := &SysExecutor{}
	se.Behaviormodel = model.NewBehaviorModel(_sim_name)
	se.dmc = NewDMC(0, definition.Infinite, "dc", "default")
	se.EXTERNAL_SRC = "SRC"
	se.EXTERNAL_DST = "DST"
	se.global_time = 0
	se.target_time = 0
	se.time_step = _time_step
	se.simulation_mode = definition.SIMULATION_IDLE
	se.sim_mode = _sim_mode
	se.waiting_obj_map = make(map[float64][]*BehaviorModelExecutor)
	se.active_obj_map = make(map[float64]*BehaviorModelExecutor)
	se.port_map = make(map[Object][]Object)
	se.Register_entity(se.dmc.executor)
	// se.min_schedule_item
	// se.output_event_queue = *deque.New()
	se.sim_init_time = time.Now()
	se.input_event_queue = input_heap{}
	heap.Init(&se.input_event_queue)
	return se
}

func (se SysExecutor) Get_global_time() float64 {
	return se.global_time
}

func (se *SysExecutor) Register_entity(sim_obj *BehaviorModelExecutor) {
	se.waiting_obj_map[sim_obj.Get_create_time()] = append(se.waiting_obj_map[sim_obj.Get_create_time()], sim_obj)
	// waiting_obj_map 에 create_time 별로 슬라이스를 만들어서 sim_obj 를 append 한다.
	fmt.Println("\n Register_entity :", time.Since(Start_time))
}

func (se *SysExecutor) Create_entity() {
	if len(se.waiting_obj_map) != 0 {
		key, value := func() (float64, []*BehaviorModelExecutor) {
			var key float64 = definition.Infinite
			for k := range se.waiting_obj_map {
				if k < key {
					key = k
				}
			}
			value := se.waiting_obj_map[key]
			return key, value
		}() //key = create_time, value = obj의 슬라이스
		if key <= se.global_time {
			for _, v := range value {
				se.active_obj_map[float64(v.sysobject.Get_obj_id())] = v
				v.Set_req_time(se.global_time, 0) //elpased ti
				se.min_schedule_item = append(se.min_schedule_item, v)
				//슬라이스를 순회하여 obj 를 active_obj_map 에 넣는다.
			}
			delete(se.waiting_obj_map, key)

			sort.Sort(schedule_item(se.min_schedule_item))

		}
	}
	fmt.Println("\n Create_entity :", time.Since(Start_time))
}

func (se *SysExecutor) Destory_entity() {

	if len(se.active_obj_map) != 0 { //active obj map 에 obj 가 있으면
		var delete_lst []*BehaviorModelExecutor
		// var port_del_lst []string
		for _, agent := range se.active_obj_map {
			if agent.Get_destruct_time() <= se.global_time { //active_obj_map을 순회하고,
				delete_lst = append(delete_lst, agent) // 이미생성된 obj 들을 delete_lst에 담고
			}
		}
		for _, agent := range delete_lst {
			delete(se.active_obj_map, float64(agent.sysobject.Get_obj_id())) // delete_lst를 순회하여 active_obj_map 에 있는 obj 를 지운다.
			var port_del_lst []Object
			for k, v := range se.port_map {
				if v[0].object == agent { //지운 obj 와 연결되어있는 port 를 port_map에서 지운다.
					port_del_lst = append(port_del_lst, k)
				}
			}
			for _, v := range port_del_lst {
				delete(se.port_map, v)
			}
			for i, v := range se.min_schedule_item {
				if v == agent {
					se.min_schedule_item = remove(se.min_schedule_item, i)
				}
			}
			//mim_schedule_item에서도 지운다.
		}
	}
	fmt.Println("\n Destory_entity :", time.Since(Start_time))
}

func (se *SysExecutor) Coupling_relation(src_obj *BehaviorModelExecutor, out_port string, dst_obj *BehaviorModelExecutor, in_port string) {

	dst := Object{dst_obj, in_port}
	b := func() bool {
		for k := range se.port_map {
			if k.object == src_obj && k.port == out_port {
				se.port_map[k] = append(se.port_map[k], dst)
				return true //port_map 에 이미있으면 추가
			}
		}
		return false
	}()
	if !b { // 없으면 새로만든다.
		src := Object{src_obj, out_port}
		se.port_map[src] = append(se.port_map[src], dst)
	}
	fmt.Println("\n Coupling_relation :", time.Since(Start_time))
}

func (se *SysExecutor) Single_output_handling(obj *BehaviorModelExecutor, msg *system.SysMessage) {

	pair := Object{obj, msg.Get_dst()}

	b := func() bool {
		for k := range se.port_map {
			if k.object == obj {
				return true
			}
		}
		return false
	}()

	if !b {
		dmc := Object{se.active_obj_map[float64(se.dmc.executor.sysobject.Get_obj_id())], "uncaught"}
		se.port_map[pair] = append(se.port_map[pair], dmc)
	}

	dst := se.port_map[pair]
	if dst == nil { //도착지가없다
		err := func() error {
			return errors.New("destination not found")
		}()
		fmt.Println(err)
	}

	for _, v := range dst {
		if v.object == nil {
			e := o_event_queue{se.global_time, msg.Retrieve()}
			se.output_event_queue = append([]*o_event_queue{&e}, se.output_event_queue...)
		} else {
			v.object.Ext_trans(v.port, msg) // msg.retrieve()
			v.object.Set_req_time(se.global_time, 0)
		}
	}

	fmt.Println("\n Single_output_handling :", time.Since(Start_time))

}

func (se *SysExecutor) output_handling(obj *BehaviorModelExecutor, msg *system.SysMessage) {

	if !(msg == nil) {
		se.Single_output_handling(obj, msg)
	}
	fmt.Println("\n output_handling :", time.Since(Start_time))
}

func (se *SysExecutor) Init_sim() {
	se.simulation_mode = definition.SIMULATION_RUNNING

	if se.active_obj_map == nil {
		se.global_time = 0
	}

	if len(se.min_schedule_item) != 0 {
		for _, obj := range se.active_obj_map {
			if obj.Time_advance() < 0 {
				err := func() error {
					return errors.New("you should give posistive real number for the deadline")
				}()
				fmt.Println(err)
			}
			obj.Set_req_time(se.global_time, 0)
			se.min_schedule_item = append(se.min_schedule_item, obj)
		}
	}
	fmt.Println("\n Init_sim :", time.Since(Start_time))
}

func (se *SysExecutor) Schedule() {
	se.Create_entity()
	se.Handle_external_input_event()

	tuple_obj := se.min_schedule_item[0]
	se.min_schedule_item = remove(se.min_schedule_item, 0)

	// before := time.Now()

	fmt.Println("global time :", se.global_time, "obj:", tuple_obj, "req_time :", tuple_obj.Get_req_time())

	const epsilon = 1e-14
	// start_time := time.Now()
	for {

		t := math.Abs(tuple_obj.Get_req_time() - se.global_time)
		if t > epsilon {
			break
		}

		msg := tuple_obj.Output()

		if msg != nil {
			se.output_handling(tuple_obj, msg)
		}
		tuple_obj.Int_trans()
		req_t := tuple_obj.Get_req_time()
		tuple_obj.Set_req_time(req_t, 0)
		se.min_schedule_item = append(se.min_schedule_item, tuple_obj)

		sort.Sort(schedule_item(se.min_schedule_item))

		tuple_obj = se.min_schedule_item[0]
		se.min_schedule_item = remove(se.min_schedule_item, 0)

		fmt.Println("obj : ", tuple_obj)
		fmt.Println("req_time :", tuple_obj.Get_req_time())
	}

	se.min_schedule_item = append([]*BehaviorModelExecutor{tuple_obj}, se.min_schedule_item...)
	fmt.Println()
	// after := time.Since(before)

	// if se.sim_mode == "REAL_TIME" {
	// 	x := se.time_step - float64(after)
	// 	if x < 0 {
	// 		time.Sleep(0)
	// 	} else {
	// 		//time.sleep(x)
	// 		time.Sleep(1 * time.Duration(x))
	// 	}

	// }
	se.global_time += se.time_step
	se.Destory_entity()
	fmt.Println("\n Schedule :", time.Since(Start_time))

}

func (se *SysExecutor) Simulate(_time float64) { //default = infinity
	fmt.Println("sumulate")
	se.target_time = se.global_time + _time
	se.Init_sim()
	for se.global_time < se.target_time {
		if se.waiting_obj_map == nil {
			item := se.min_schedule_item[0]
			se.min_schedule_item = remove(se.min_schedule_item, 0)
			if item.Get_req_time() == definition.Infinite && se.sim_mode == "VIRTURE_TIME" {
				se.simulation_mode = definition.SIMULATION_TERMINATED
				break
			}
			se.min_schedule_item = append([]*BehaviorModelExecutor{item}, se.min_schedule_item...)
		}
		se.Schedule()
	}

}

func (se *SysExecutor) Simulation_stop() {
	se.global_time = 0
	se.target_time = 0
	se.time_step = 1
	se.waiting_obj_map = make(map[float64][]*BehaviorModelExecutor)
	se.active_obj_map = make(map[float64]*BehaviorModelExecutor)
	se.port_map = make(map[Object][]Object)
	// se.min_schedule_item = *deque.New()
	se.sim_init_time = time.Now()
	se.dmc = NewDMC(0, definition.Infinite, "dc", "default")
	se.Register_entity(se.dmc.executor)
}

func (se *SysExecutor) Insert_external_event(_port string, _msg interface{}, scheduled_time float64) {
	sm := system.NewSysMessage("SRC", _port)
	sm.Insert(_msg)

	_, b := Slice_Find_string(se.Behaviormodel.CoreModel.Intput_ports, _port)
	if b {
		//lock.acquire
		eq := i_event_queue{scheduled_time + se.global_time, sm}
		heap.Push(&se.input_event_queue, eq)
		//lock.release()
	} else {
		print("[ERROR][INSERT_EXTERNAL_EVNT] Port Not Found")
	}
	fmt.Println("\n Insert_external_event :", time.Since(Start_time))

}

func (se *SysExecutor) Insert_custom_external_event(_port string, _bodylist []interface{}, scheduled_time float64) {
	sm := system.NewSysMessage("SRC", _port)
	sm.Extend(_bodylist)
	_, b := Slice_Find_string(se.Behaviormodel.CoreModel.Intput_ports, _port)
	if b {
		//lock.acquire
		eq := i_event_queue{scheduled_time + se.global_time, sm}
		heap.Push(&se.input_event_queue, eq)
		//lock.release()
	} else {
		fmt.Printf("[ERROR][INSERT_EXTERNAL_EVNT] Port Not Found")
	}
	fmt.Println("\n Insert_custom_external_event :", time.Since(Start_time))
}

func (se *SysExecutor) Get_generated_event() []*o_event_queue {
	return se.output_event_queue
}

func (se *SysExecutor) Handle_external_input_event() {

	var event_list []i_event_queue
	for _, ev := range se.input_event_queue {
		if ev.time <= se.global_time {
			event_list = append(event_list, ev)
		}
	}
	for _, event := range event_list {

		se.output_handling(nil, event.msg)
		heap.Pop(&se.input_event_queue)

	}

	sort.Sort(schedule_item(se.min_schedule_item))
	fmt.Println("\n Handle_external_input_event :", time.Since(Start_time))
}

func (se *SysExecutor) Handle_external_output_event() deque.Deque {
	fmt.Println("handle_external_output_event")
	var event_lists deque.Deque
	err := deepcopy.Copy(event_lists, se.output_event_queue)
	if err != nil {
		err := func() error {
			return errors.New("can't Handle_external_output_event")
		}()
		fmt.Println(err)
	}
	se.output_event_queue = se.output_event_queue[:0]
	return event_lists
}

func (se *SysExecutor) Is_terminated() interface{} {
	return se.simulation_mode == definition.SIMULATION_TERMINATED
}
