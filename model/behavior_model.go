package model

import (
	"evsim_golang/definition"
)

type Behaviormodel struct {
	States    map[string]float64
	CoreModel *definition.CoreModel
}

func (b *Behaviormodel) Insert_state(name string, deadline float64) { //deadline 디폴트값 = 0
	b.States[name] = deadline
}

func (b *Behaviormodel) Update_state(name string, deadline float64) { //deadline 디폴트값 = 0
	b.States[name] = deadline
}

func NewBehaviorModel(name string) *Behaviormodel {
	b := Behaviormodel{}
	b.States = make(map[string]float64)
	b.CoreModel = definition.NewCoreModel(name, definition.BEHAVIORAL)
	return &b
}
