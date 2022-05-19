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

var start_time time.Time

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
	min_schedule_item  deque.Deque
	input_event_queue  input_heap
	output_event_queue deque.Deque
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
	se.min_schedule_item = *deque.New()
	se.output_event_queue = *deque.New()
	se.sim_init_time = time.Now()
	se.input_event_queue = input_heap{}
	heap.Init(&se.input_event_queue)
	return se
}

func (se SysExecutor) Get_global_time() float64 {
	return se.global_time
}

func (se *SysExecutor) Register_entity(sim_obj *BehaviorModelExecutor) {
	fmt.Println("Register_entity")
	se.waiting_obj_map[sim_obj.Get_create_time()] = append(se.waiting_obj_map[sim_obj.Get_create_time()], sim_obj)
	// waiting_obj_map 에 create_time 별로 슬라이스를 만들어서 sim_obj 를 append 한다.
}

func (se *SysExecutor) Create_entity() {
	fmt.Println("create_entity")
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
				se.min_schedule_item.PushBack(v)
				//슬라이스를 순회하여 obj 를 active_obj_map 에 넣는다.
			}
			delete(se.waiting_obj_map, key)

			var lst []*BehaviorModelExecutor
			len := se.min_schedule_item.Len()
			for i := 0; i < len; i++ {
				lst = append(lst, se.min_schedule_item.PopFront().(*BehaviorModelExecutor))
			}
			sort.Sort(schedule_item(lst))
			for i := 0; i < len; i++ {
				se.min_schedule_item.PushBack(lst[i])
			}

			t := time.Now()
			time1 := t.Sub(start_time)
			fmt.Println("\n create_entity_time :", time1)
		}
	}
}

func (se *SysExecutor) Destory_entity() {
	fmt.Println("Destory_entity")
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
			i := se.min_schedule_item.Index(func(i interface{}) bool {
				if i == agent {
					return true
				} else {
					return false
				}
			})
			se.min_schedule_item.Remove(i)
			//mim_schedule_item에서도 지운다.
		}
	}
}

func (se *SysExecutor) Coupling_relation(src_obj *BehaviorModelExecutor, out_port string, dst_obj *BehaviorModelExecutor, in_port string) {
	fmt.Println("coupling_relation")
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
}

func (se *SysExecutor) Single_output_handling(obj *BehaviorModelExecutor, msg *system.SysMessage) {
	fmt.Println("single_Output_handling")
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
			se.output_event_queue.PushFront(e)
		} else {
			v.object.Ext_trans(v.port, msg) // msg.retrieve()
			v.object.Set_req_time(se.global_time, 0)
		}
	}

}

func (se *SysExecutor) output_handling(obj *BehaviorModelExecutor, msg *system.SysMessage) {
	fmt.Println("output_handling")
	if !(msg == nil) {
		se.Single_output_handling(obj, msg)
	}
}

func (se *SysExecutor) Init_sim() {
	fmt.Println("Init_sim")
	se.simulation_mode = definition.SIMULATION_RUNNING

	if se.active_obj_map == nil {
		se.global_time = 0
	}

	if se.min_schedule_item.Cap() != 0 {
		for _, obj := range se.active_obj_map {
			if obj.Time_advance() < 0 {
				err := func() error {
					return errors.New("you should give posistive real number for the deadline")
				}()
				fmt.Println(err)
			}
			obj.Set_req_time(se.global_time, 0)
			se.min_schedule_item.PushBack(obj)
		}
	}
}

func (se *SysExecutor) Schedule() {
	fmt.Println("schedule")
	se.Create_entity()
	se.Handle_external_input_event()

	tuple_obj := se.min_schedule_item.PopFront().(*BehaviorModelExecutor)
	// before := time.Now()

	fmt.Println("global time :", se.global_time, "obj:", tuple_obj, "req_time :", tuple_obj.Get_req_time())

	const epsilon = 1e-14

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
		se.min_schedule_item.PushBack(tuple_obj)

		var lst []*BehaviorModelExecutor
		len := se.min_schedule_item.Len()
		for i := 0; i < len; i++ {
			lst = append(lst, se.min_schedule_item.PopFront().(*BehaviorModelExecutor))
		}
		sort.Sort(schedule_item(lst))
		for i := 0; i < len; i++ {
			se.min_schedule_item.PushBack(lst[i])
		}

		tt := time.Now()
		time2 := tt.Sub(start_time)
		fmt.Println("\n schedule_time :", time2)

		tuple_obj = se.min_schedule_item.PopFront().(*BehaviorModelExecutor)
		fmt.Println("obj : ", tuple_obj)
		fmt.Println("req_time :", tuple_obj.Get_req_time())
	}

	se.min_schedule_item.PushFront(tuple_obj)
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

}

func (se *SysExecutor) Simulate(_time float64) { //default = infinity
	fmt.Println("sumulate")
	se.target_time = se.global_time + _time
	se.Init_sim()
	start_time = time.Now()
	for se.global_time < se.target_time {
		if se.waiting_obj_map == nil {
			item := se.min_schedule_item.PopFront().(*BehaviorModelExecutor)
			if item.Get_req_time() == definition.Infinite && se.sim_mode == "VIRTURE_TIME" {
				se.simulation_mode = definition.SIMULATION_TERMINATED
				break
			}
			se.min_schedule_item.PushFront(item)
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
	se.min_schedule_item = *deque.New()
	se.sim_init_time = time.Now()
	se.dmc = NewDMC(0, definition.Infinite, "dc", "default")
	se.Register_entity(se.dmc.executor)
}

func (se *SysExecutor) Insert_external_event(_port string, _msg interface{}, scheduled_time float64) {
	fmt.Println("insert_external_event")
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

}

func (se *SysExecutor) Insert_custom_external_event(_port string, _bodylist []interface{}, scheduled_time float64) {
	fmt.Println("insert_custom_external_event")
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
}

func (se *SysExecutor) Get_generated_event() deque.Deque {
	return se.output_event_queue
}

func (se *SysExecutor) Handle_external_input_event() {
	fmt.Println("handle_external_input_event")
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

	var lst []*BehaviorModelExecutor
	len := se.min_schedule_item.Len()
	for i := 0; i < len; i++ {
		lst = append(lst, se.min_schedule_item.PopFront().(*BehaviorModelExecutor))
	}
	sort.Sort(schedule_item(lst))
	for i := 0; i < len; i++ {
		se.min_schedule_item.PushBack(lst[i])
	}

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
	se.output_event_queue.Clear()
	return event_lists
}

func (se *SysExecutor) Is_terminated() interface{} {
	return se.simulation_mode == definition.SIMULATION_TERMINATED
}
