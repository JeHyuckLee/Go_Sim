package main

import "evsim_golang/executor"

type SystemSimulator struct {
	engine map[string]*executor.SysExecutor
}

func (ss *SystemSimulator) Register_engine(sim_name, sim_mode string, time_step int) { //time_step =1 (default)
	ss.engine[sim_name] = executor.NewSysExecutor(time_step, sim_name, sim_mode)
}

func (ss SystemSimulator) Get_engine_map() map[string]*executor.SysExecutor {
	return ss.engine
}
func (ss SystemSimulator) Get_engine(sim_name string) *executor.SysExecutor {
	return ss.engine[sim_name]
}

func (ss SystemSimulator) Is_terminated(sim_name string) interface{} {
	return ss.engine[sim_name].Is_terminated()
}

func (ss SystemSimulator) Set_learning_module(sim_name string, learn_module interface{}) {
	ss.engine[sim_name].Set_learning_module(learn_module)
}

func (ss SystemSimulator) Get_learning_module(sim_name string) interface{} {
	return ss.engine[sim_name].Get_learning_module()
}

func (ss SystemSimulator) Exec_simulation_instance(instance_path interface{}) {

}

func NewSysSimulator() *SystemSimulator {
	ss := &SystemSimulator{}
	return ss
}
