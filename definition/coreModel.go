package definition

type CoreModel struct {
	_type         int
	_name         string
	_intput_ports []string
	_output_ports []string
}

func (c *CoreModel) Set_name(name string) {
	name = c._name
}

func (c *CoreModel) Get_name() string {
	return c._name
}
func (c *CoreModel) Insert_input_port(port string) {
	c._intput_ports = append(c._intput_ports, port)
}
func (c *CoreModel) Retrieve_input_port() []string {
	return c._intput_ports
}

func (c *CoreModel) Insert_output_port(port string) {
	c._output_ports = append(c._output_ports, port)
}

func (c *CoreModel) Retrieve_output_port() []string {
	return c._output_ports
}

func (c *CoreModel) Get_type() int {
	return c._type
}

func NewCoreModel(_name string, _type int) *CoreModel {
	c := CoreModel{}
	c._name = _name
	c._type = _type
	return &c
}
