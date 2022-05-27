package definition

type CoreModel struct {
	_type         int
	Name          string
	Intput_ports  []string
	_output_ports []string
}

func (c *CoreModel) Set_name(name string) {
	c.Name = name
}

func (c *CoreModel) Get_name() string {
	return c.Name
}
func (c *CoreModel) Insert_input_port(port string) {
	c.Intput_ports = append(c.Intput_ports, port)
}
func (c *CoreModel) RetrieveInput_port() []string {
	return c.Intput_ports
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
	c.Name = _name
	c._type = _type
	return &c
}
