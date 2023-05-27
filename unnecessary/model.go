package unnecessary

import (
	"fmt"
	"log"
	"reflect"
)

// type Model interface{}

// type DynamicModel func() Model

// func GetModelAsText(model Model) string {
// 	// var value Model
// 	// log.Printf("getModelAsText - %t", model)
// 	switch v := model.(type) {
// 	// case Model2:
// 	// 	log.Printf("Model2!!")
// 	// 	return GetModelAsText(v.String())
// 	// case DynamicModel:
// 	// 	// log.Printf("DynamicModel")
// 	// 	return GetModelAsText(v())
// 	case func() interface{}:
// 		// log.Printf("func() interface {}")
// 		return GetModelAsText(v())
// 	case string:
// 		// log.Printf("string")
// 		return v
// 	default:
// 		// log.Printf("other")
// 		return fmt.Sprintf("%v", v)
// 	}
// }

func ValueToString(value interface{}) string {
	// log.Printf("1 - %v", reflect.ValueOf(model.Value).Type().Kind() == reflect.Pointer)
	rvalue := reflect.ValueOf(value)
	if rvalue.Type().Kind() == reflect.Pointer {
		value = reflect.Indirect(rvalue).Interface()
	}
	switch v := value.(type) {
	case string:
		return v
	default:
		return fmt.Sprintf("%v", value)
	}
}

func ToListOfModel(values ...interface{}) []Model {
	res := make([]Model, len(values))
	for i, value := range values {
		switch v := value.(type) {
		case Model:
			res[i] = v
		default:
			res[i] = NewGenericModel(value)
		}
	}
	return res
}

type Model interface {
	String() string
	List() []Model
}

type GenericModel struct {
	Value interface{}
}

func NewGenericModel(value interface{}) Model {
	return &GenericModel{Value: value}
}

func (model *GenericModel) String() string {
	return ValueToString(model.Value)
}

func (model *GenericModel) List() []Model {
	panic(fmt.Errorf("the operation 'List()' is not support by %T", model))
}

func (model *GenericModel) SetValue(value interface{}) {
	model.Value = value
}

type StringModel struct {
	Value string
}

func NewStringModel(value string) Model {
	return &StringModel{Value: value}
}

func (model *StringModel) String() string {
	return model.Value
}

func (model *StringModel) List() []Model {
	panic(fmt.Errorf("the operation 'List()' is not support by %T", model))
}

func (model *StringModel) SetValue(value string) {
	model.Value = value
}

type GenericListModel struct {
	Value []Model
}

func NewGenericListModel(values ...interface{}) Model {
	return &GenericListModel{Value: ToListOfModel(values...)}
}

func (model *GenericListModel) String() string {
	panic(fmt.Errorf("the operation 'String()' is not support by %T", model))
}

func (model *GenericListModel) List() []Model {
	return model.Value
}

func (model *GenericListModel) SetValue(values ...interface{}) {
	model.Value = ToListOfModel(values...)
}

type DynamicModelCallback func() interface{}

type DynamicModel struct {
	Callback DynamicModelCallback
}

func NewDynamicModel(callback DynamicModelCallback) Model {
	return &DynamicModel{Callback: callback}
}

func (model *DynamicModel) String() string {
	value := model.Callback()
	return ValueToString(value)
}

func (model *DynamicModel) List() []Model {
	panic(fmt.Errorf("the operation 'List()' is not support by %T", model))
}

func PlayWithModels() {
	model0 := GenericModel{Value: nil}
	log.Printf("model0 -> %s", &model0)

	model1 := GenericModel{Value: 1}
	log.Printf("model1 -> %s", &model1)

	model2 := NewGenericModel("olala")
	log.Printf("model2 -> %s", model2)

	model3 := GenericListModel{Value: []Model{&GenericModel{Value: 1}}}
	log.Printf("model2 -> %v", model3.List())

	model4 := NewGenericListModel(1, 2, "a", NewGenericModel(12))
	log.Printf("model2 -> %v", model4.List())

	value5 := 123
	model5 := GenericModel{Value: &value5}
	log.Printf("model5 -> %s", model5.String())
}
