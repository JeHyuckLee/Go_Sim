package main

import (
	"evsim_golang/executor"
	"evsim_golang/system"
)

type move struct {
	executor *executor.BehaviorModelExecutor
	msg_list []interface{}
}

func (m *move) Int_trans() {

}

func (m *move) Ext_trans(port string, msg *system.SysMessage) {

}

func (m *move) Output() *system.SysMessage {

}
