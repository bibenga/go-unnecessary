package unnecessary

import (
	"fmt"
	"reflect"
)

func ValueToString(value interface{}) string {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Pointer:
		rvalue := reflect.ValueOf(value)
		value = reflect.Indirect(rvalue).Interface()
	}
	switch v := value.(type) {
	case string:
		return v
	default:
		return fmt.Sprintf("%v", value)
	}
}

func ValueIsSlice(value interface{}) bool {
	return reflect.TypeOf(value).Kind() == reflect.Slice
}

func ValueToSlice(value interface{}) ([]interface{}, error) {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Pointer:
		rvalue := reflect.ValueOf(value)
		value = reflect.Indirect(rvalue).Interface()
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice:
		rvalue := reflect.ValueOf(value)
		out := make([]interface{}, rvalue.Len())
		for i := 0; i < rvalue.Len(); i++ {
			out[i] = rvalue.Index(i).Interface()
		}
		return out, nil
	default:
		return nil, fmt.Errorf("value is not slice")
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
