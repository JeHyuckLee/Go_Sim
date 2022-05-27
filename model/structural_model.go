package model

import "evsim_golang/definition"

type StructuralModel struct {
	name      string
	_models   []string
	CoreModel *definition.CoreModel
}

func NewStructuralModel(name string) *StructuralModel {
	str := StructuralModel{}
	str.CoreModel = definition.NewCoreModel(name, definition.BEHAVIORAL)
	return &str
}

func (str *StructuralModel) Insert_model(model string) {
	str._models = append(str._models, model)
}

func (str *StructuralModel) Retrieve_models() []string {
	return str._models
}
