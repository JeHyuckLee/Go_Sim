package behaviormodel

import (
	"evsim_golang/definition"
	"strconv"
)

type Behaviormodel struct {
	_states                       map[string]float64
	external_transition_map_tuple []string
	external_transition_map_state []string
	internal_transition_map_tuple []string
	internal_transition_map_state []string
	coreModel                     *definition.CoreModel
}

func (b *Behaviormodel) Insert_state(name, deadline string) {
	num, _ := strconv.ParseFloat(deadline, 64)
	b._states[name] = num

}

func (b *Behaviormodel) Update_state(name, deadline string) {
	num, _ := strconv.ParseFloat(deadline, 64)
	b._states[name] = num

}

func NewBehaviorModel(name string) *Behaviormodel {
	b := Behaviormodel{}
	b.coreModel = definition.NewCoreModel(name, 0)
	return &b
}
