package executor

import (
	"evsim_golang/definition"
	"evsim_golang/model"
	"evsim_golang/my"
	"evsim_golang/system"
	"fmt"
	"time"

	"github.com/gammazero/deque"
	"gopkg.in/getlantern/deepcopy.v1"
)

type SysExecutor struct {
	sysObject     *system.SysObject
	behaviormodel *model.Behaviormodel
	dmc           *DefaultMessageCatcher

	global_time        float64
	target_time        float64
	time_step          int
	EXTERNAL_SRC       string
	EXTERNAL_DST       string
	simulation_mode    int
	min_schedule_item  deque.Deque
	input_event_queue  []string
	output_event_queue deque.Deque
	sim_mode           string
	waiting_obj_map    map[float64][]*BehaviorModelExecutor
	active_obj_map     map[float64]*BehaviorModelExecutor
	learn_module       interface{}
	port_map           map[interface{}]interface{}
	sim_init_time      time.Time
}

func NewSysExecutor(_time_step int, _sim_name, _sim_mode string) *SysExecutor {
	se := SysExecutor{}
	se.behaviormodel = model.NewBehaviorModel(_sim_name)
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
	se.Register_entity(se.dmc.executor)
	se.min_schedule_item = *deque.New()
	se.output_event_queue = *deque.New()
	se.sim_init_time = time.Now()
}

func (se SysExecutor) Get_global_time() float64 {
	return se.global_time
}

func (se *SysExecutor) Register_entity(sim_obj *BehaviorModelExecutor) {
	se.waiting_obj_map[sim_obj.Get_create_time()] = append(se.waiting_obj_map[sim_obj.Get_create_time()], sim_obj)
	// waiting_obj_map 에 create_time 별로 슬라이스를 만들어서 sim_obj 를 append 한다.
}

func (se *SysExecutor) Create_entity() {
	if len(se.waiting_obj_map) != 0 {
		key, value := func() (float64, []*BehaviorModelExecutor) {
			var key float64 = 0
			for k, _ := range se.waiting_obj_map {
				if k < key {
					key = k
				}
			}
			value := se.waiting_obj_map[key]
			return key, value
		}() //key = create_time, value = obj의 슬라이스
		for _, v := range value {
			se.active_obj_map[float64(v.sysobject.Get_obj_id())] = v
			v.Set_req_time(se.global_time, 0) //elpased ti
			se.min_schedule_item.PushFront(v)
			//슬라이스를 순회하여 obj 를 active_obj_map 에 넣는다.
		}
		delete(se.waiting_obj_map, key)
		// se.min_schedule_itme 정렬
	}
}

func (se *SysExecutor) Destory_entity() {
	if len(se.active_obj_map) != 0 {
		var delete_lst []*BehaviorModelExecutor
		for _, agent := range se.active_obj_map {
			if agent.Get_create_time() <= se.global_time { // int? float??
				delete_lst = append(delete_lst, agent)
			}
		}
		for _, v := range delete_lst {
			delete(se.active_obj_map, float64(v.sysobject.Get_obj_id()))
			port_del_lst
		}

		// for _, v := range delete_lst {
		// delete(se.active_obj_map, v.executor.sysobject.Get_obj_id())
		// var port_del_lst []
		// 	for key, value in self.port_map.items():
		//             #print(value)
		//             if value:
		//                 if value[0][0] is agent:
		//                     port_del_lst.append(key)

		//         for key in port_del_lst:
		//             del (self.port_map[key])
		//         self.min_schedule_item.remove(agent)
		// }

	}
}

func (se *SysExecutor) Coupling_relation(src_obj, out_port, dst_obj, in_port interface{}) {

	// for k,v := range se.port_map{
	// 	if k == src_obj && v == dst_obj{
	// 		se.port_map[(src_obj,dst_obj)].append((dst_obj,in_port))
	// 	}else{
	// 		se.port_map[(src_obj,out_port)] = [(dst_obj,in_port)]
	// 	}
	// }

	// if (src_obj, out_port) in self.port_map:
	//         self.port_map[(src_obj, out_port)].append((dst_obj, in_port))
	//     else:
	//         self.port_map[(src_obj, out_port)] = [(dst_obj, in_port)]
}

func (se *SysExecutor) _Coupling_relation(src, dst interface{}) {
	_, bool := my.Map_Find(se.port_map, src)
	if bool == true {
		se.port_map[src] = dst
	} else {
		se.port_map[src] = dst
	}
}

func (se *SysExecutor) Single_output_handling(obj, msg interface{}) {
	p := my.Pair{obj, msg} // msg[1].get_dst()

	_, bool := my.Map_Find(se.port_map, p)
	if bool == false {
		// se.port_map[p] = [(se.active_obj_map[se.dmc.executor.sysobject.Get_obj_id()],"uncaught")]
	}

	for k, v := range se.port_map {
		if k == p {
			destination := v
			if destination == nil {
				fmt.Println("Destination Not Found")
				fmt.Println(se.port_map)
				// raise AssertionError
			}
		}
	}
	// if destination[0] is None:
	//             self.output_event_queue.append(
	//                 (self.global_time, msg[1].retrieve()))
	// else:
	// destination[0].ext_trans(destination[1], msg[1])
	// while self.thread_flag:
	//             time.sleep(0.001)

	//         self.thread_flag = True
	// destination[0].set_req_time(self.global_time)
	// self.thread_flag = False
}

func (se *SysExecutor) output_handling(obj, msg interface{}) {
	if !(msg == nil) {
		// if type(msg) == list:
		//         for ith_msg in msg:
		//             self.single_output_handling(obj, copy.deepcopy(ith_msg))
		//     else:
		//         self.single_output_handling(obj, msg)
	}
}

func (se *SysExecutor) Flattening(_model, _del_model, _del_coupling interface{}) {

}

func (se *SysExecutor) Init_sim() {
	se.simulation_mode = definition.SIMULATION_RUNNING
	var _del_model []*BehaviorModelExecutor
	// var _del_coupling []
	for _, model_list := range se.waiting_obj_map {
		for modle := range model_list {

		}
	}
}

func (se *SysExecutor) Schedule() {
	se.Create_entity()
	se.Handle_external_input_event()

	tuple_obj := se.min_schedule_item.PopFront()

	// before = time.perf_counter()  # TODO: consider decorator

	//     while math.isclose(tuple_obj.get_req_time(),
	//                        self.global_time,
	//                        rel_tol=1e-9):
	//         msg = tuple_obj.output()
	//         if msg is not None:
	//             self.output_handling(tuple_obj, (self.global_time, msg))

	//         # Sender Scheduling
	//         tuple_obj.int_trans()
	//         req_t = tuple_obj.get_req_time()

	//         tuple_obj.set_req_time(req_t)
	//         self.min_schedule_item.append(tuple_obj)

	//         self.min_schedule_item = deque(
	//             sorted(self.min_schedule_item,
	//                    key=lambda bm: bm.get_req_time()))

	//         tuple_obj = self.min_schedule_item.popleft()

	//     self.min_schedule_item.appendleft(tuple_obj)

	//     after = time.perf_counter()
	//     if self.sim_mode == "REAL_TIME":
	//         time.sleep((lambda x: x if x > 0 else 0)(float(self.time_step) -
	//                                                  float(after - before)))

	se.global_time += float64(se.time_step)
	se.Destory_entity()

}
func (se *SysExecutor) Simulate(_time float64) { //default = infinity
	se.target_time = se.global_time + _time
	se.Init_sim()

	for se.global_time < se.target_time {
		// if not self.waiting_obj_map:
		// if self.min_schedule_item[0].get_req_time(
		// ) == Infinite and self.sim_mode == 'VIRTUAL_TIME':
		// 	self.simulation_mode = SimulationMode.SIMULATION_TERMINATED

		se.Schedule()
	}

}
func (se *SysExecutor) Simulation_stop() {
	se.global_time = 0
	se.target_time = 0
	se.time_step = 1
	for k := range se.waiting_obj_map {
		delete(se.waiting_obj_map, k)
	}
	for k := range se.active_obj_map {
		delete(se.active_obj_map, k)
	}
	for k := range se.port_map {
		delete(se.port_map, k)
	}
	se.min_schedule_item = *deque.New()
	se.sim_init_time = time.Now()
	se.dmc = NewDMC(0, definition.Infinite, "dc", "default")
	se.Register_entity(se.dmc.executor)

}

func (se *SysExecutor) Insert_external_event(_port, _msg string, scheduled_time int) {
	sm := system.NewSysMessage("SRC", _port)
	sm.Insert(_msg)
	_, bool := my.Slice_Find_string(se.behaviormodel.CoreModel.Intput_ports, _port)
	if bool == true {
		// self.lock.acquire()
		// heapq.heappush(self.input_event_queue,
		//                (scheduled_time + self.global_time, sm))
		// self.lock.release()
	} else {
		print("[ERROR][INSERT_EXTERNAL_EVNT] Port Not Found")
	}

}

func (se *SysExecutor) Insert_custom_external_event(_port string, _bodylist []string, scheduled_time interface{}) {
	sm := system.NewSysMessage("SRC", _port)
	sm.Extend(_bodylist)
	_, bool := my.Slice_Find_string(se.behaviormodel.CoreModel.Intput_ports, _port)
	if bool == true {
		// self.lock.acquire()
		// heapq.heappush(self.input_event_queue,
		//                (scheduled_time + self.global_time, sm))
		// self.lock.release()
	} else {
		fmt.Printf("[ERROR][INSERT_EXTERNAL_EVNT] Port Not Found")
	}
}

func (se *SysExecutor) Get_generated_event() deque.Deque {
	return se.output_event_queue
}

func (se *SysExecutor) Handle_external_input_event() {
	// event_list = [
	//         ev for ev in self.input_event_queue if ev[0] <= self.global_time
	//     ]
	//     #print(event_list)
	//     for event in event_list:
	//         self.output_handling(None, event)
	//         self.lock.acquire()
	//         heapq.heappop(self.input_event_queue)
	//         self.lock.release()

	//     self.min_schedule_item = deque(
	//         sorted(self.min_schedule_item, key=lambda bm: bm.get_req_time()))
	//     pass
}

func (se *SysExecutor) Handle_external_output_event() deque.Deque {
	var event_lists deque.Deque
	err := deepcopy.Copy(event_lists, se.output_event_queue)
	if err != nil {

	}
	se.output_event_queue.Clear()
	return event_lists
}

func (se *SysExecutor) Is_terminated() interface{} {
	return se.simulation_mode == definition.SIMULATION_TERMINATED
}

func (se SysExecutor) Set_learning_module(learn_module interface{}) {
	se.learn_module = learn_module
}

func (se *SysExecutor) Get_learning_module() interface{} {
	return se.learn_module
}
