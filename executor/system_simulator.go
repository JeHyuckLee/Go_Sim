package executor

type SystemSimulator struct {
	Engine map[string]*SysExecutor
}

func (ss *SystemSimulator) Register_engine(sim_name, sim_mode string, time_step float64) { //time_step =1 (default)
	ss.Engine[sim_name] = NewSysExecutor(time_step, sim_name, sim_mode)
}

func (ss SystemSimulator) Get_engine_map() map[string]*SysExecutor {
	return ss.Engine
}
func (ss SystemSimulator) Get_engine(sim_name string) *SysExecutor {
	return ss.Engine[sim_name]
}

func (ss SystemSimulator) Is_terminated(sim_name string) interface{} {
	return ss.Engine[sim_name].Is_terminated()
}

func (ss SystemSimulator) Exec_simulation_instance(instance_path interface{}) {

}

func NewSysSimulator() *SystemSimulator {
	ss := &SystemSimulator{}
	ss.Engine = make(map[string]*SysExecutor)
	return ss
}
