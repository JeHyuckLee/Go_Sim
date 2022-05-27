package system

type SysMessage struct {
	sysobject *SysObject
	_src      string
	_dst      string
	_msg_time int
	_msg_list []interface{}
}

func (sm *SysMessage) Insert(msg interface{}) {
	sm._msg_list = append(sm._msg_list, msg)
}

func (sm *SysMessage) Extend(_list []interface{}) {
	sm._msg_list = append(sm._msg_list, _list...)
}
func (sm *SysMessage) Retrieve() []interface{} {
	return sm._msg_list
}
func (sm *SysMessage) Get_src() string {
	return sm._src
}
func (sm *SysMessage) Get_dst() string {
	return sm._dst
}
func (sm *SysMessage) Set_msg_time(t int) {
	sm._msg_time = t
}
func (sm *SysMessage) Get_msg_time() int {
	return sm._msg_time
}

func NewSysMessage(src_name, dst_name string) *SysMessage {
	sm := SysMessage{}
	sm.sysobject = NewSysObject()
	sm._src = src_name
	sm._dst = dst_name
	sm._msg_time = -1

	return &sm
}
