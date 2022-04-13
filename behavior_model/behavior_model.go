package behaviormodel

import (
	"evsim_golang/definition"
	"math"
)

type Behaviormodel struct {
	_states                       map[string]float64
	external_transition_map_tuple []string
	external_transition_map_state []string
	internal_transition_map_tuple []string
	internal_transition_map_state []string
	coreModel                     *definition.CoreModel
}

func (b *Behaviormodel) Insert_state(name string, deadline float64) { //deadline 디폴트값 = 0
	if deadline == 0 {
		deadline = math.Inf(1)
	}
	b._states[name] = deadline
}

func (b *Behaviormodel) Update_state(name string, deadline float64) { //deadline 디폴트값 = 0
	if deadline == 0 {
		deadline = math.Inf(1)
	}
	b._states[name] = deadline
}

func NewBehaviorModel(name string) *Behaviormodel {
	b := Behaviormodel{}
	b.coreModel = definition.NewCoreModel(name, 0)
	return &b
}
