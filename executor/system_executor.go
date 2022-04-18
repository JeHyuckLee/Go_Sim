package executor

import (
	"evsim_golang/definition"
	"evsim_golang/model"
	"math"
	"sort"
	"github.com/gammazero/deque"
)

type SysExecutor struct {
	sysObject          *SysObject
	behaviormodel      *model.Behaviormodel
	dmc                *DefaultMessageCatcher
	global_time        int
	target_time        int
	time_step          int
	EXTERNAL_SRC       string
	EXTERNAL_DST       string
	simulation_mode    int
	min_schedule_item  deque.Deque
	input_event_queue  []string
	output_event_queue deque.Deque
	sim_mode           string
	waiting_obj_map map[int]*DefaultMessageCatcher
	active_obj_map map[int]*DefaultMessageCatcher
}

func NewSysExecutor(_time_step int, _sim_name, _sim_mode string) *SysExecutor {
	se := SysExecutor{}
	se.behaviormodel = model.NewBehaviorModel(_sim_name)
	se.dmc = message.NewDMC(0, definition.Infinite, "dc", "default")
	se.Register_entity(se.dmc)
	se.EXTERNAL_SRC = "SRC"
	se.EXTERNAL_DST = "DST"
	se.global_time = 0
	se.target_time = 0
	se.time_step = _time_step
	se.simulation_mode = definition.SIMULATION_IDLE
	se.sim_mode = _sim_mode
	se.waiting_obj_map = make(map[int]*DefaultMessageCatcher)
	se.active_obj_map = make(map[int]*DefaultMessageCatcher)
}

func (se SysExecutor) Get_global_time() int {0
	return se.global_time
}

func min(numbers map[int]*DefaultMessageCatcher) int {
    var minNumber int
    for n := range numbers {
        minNumber = n
        break
    }
    for n := range numbers {
        if n < maxNumber {
            maxNumber = n
        }
    }
    return maxNumber
}

func (se *SysExecutor) Register_entity(sim_obj *DefaultMessageCatcher) {

}

func (se *SysExecutor) Create_entity() {
	if len(se.waiting_obj_map)!=0{
		key := min(se.waiting_obj_map)
		if key <= se.global_time{
			lst := se.waiting_obj_map[key]
			for i,obj := range lst {
				se.active_obj_map[obj.Get_obj_id()] = obj
				obj.Set_req_time(se.global_time)
				se.min_schedule_item.PushBack(obj)
			}
			delete(se.waiting_obj_map,key)
			// self.min_schedule_item = deque(
			// 	sorted(self.min_schedule_item,
			// 		   key=lambda bm: bm.get_req_time()))
			}
		}
	}
}
func (se *SysExecutor) Destory_entity() {
	if len(se.active_obj_map)!=0{
		// delete_lst = []
        //     for agent_name, agent in self.active_obj_map.items():
        //         if agent.get_destruct_time() <= self.global_time:
        //             delete_lst.append(agent)

        //     for agent in delete_lst:
        //         #print("global:",self.global_time," del agent:", agent.get_name())
        //         del (self.active_obj_map[agent.get_obj_id()])

        //         port_del_lst = []
        //         for key, value in self.port_map.items():
        //             #print(value)
        //             if value:
        //                 if value[0][0] is agent:
        //                     port_del_lst.append(key)

        //         for key in port_del_lst:
        //             del (self.port_map[key])
        //         self.min_schedule_item.remove(agent)
	}
}
func (se *SysExecutor) Coupling_relation(src_obj, out_port, dst_obj, in_port) {
	// if (src_obj, out_port) in self.port_map:
    //         self.port_map[(src_obj, out_port)].append((dst_obj, in_port))
    //     else:
    //         self.port_map[(src_obj, out_port)] = [(dst_obj, in_port)]
    //         # self.port_map_wName.append((src_obj.get_name(), out_port, dst_obj.get_name(), in_port))
}

func (se *SysExecutor) _Coupling_relation(src, dst) {
	// if src in self.port_map:
	// self.port_map[src].append(dst)
	// else:
	// 	self.port_map[src] = [dst]
}

func (se *SysExecutor) Single_output_handling(obj, msg) {

}

func (se *SysExecutor) output_handling(obj, msg) {

}

func (se *SysExecutor) Flattening(_model, _del_model, _del_coupling) {

}

func (se *SysExecutor) Init_sim() {

}

func (se *SysExecutor) Schedule() {

}

func (se *SysExecutor) Simulation_stop() {

}

func (se *SysExecutor) Insert_external_event(_port, _msg, scheduled_time int) {

}

func (se *SysExecutor) Insert_custom_external_event(_port, _bodylist, scheduled_time) {

}

func (se *SysExecutor) Get_generated_event() {

}

func (se *SysExecutor) Handle_external_input_event() {

}

func (se *SysExecutor) Handle_external_output_event() {

}

func (se *SysExecutor) Is_terminated() {

}

func (se *SysExecutor) Set_learning_module() {

}

func (se *SysExecutor) Get_learning_module() {

}
