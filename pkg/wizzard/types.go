package wizzard

import (
	"reflect"
)

type Struct struct {
	Type    reflect.Type
	Name    string
	Fields  []Field
	Methods []Method
}

type Field struct {
	Type  reflect.Type
	Name  string
	Index int
}

type Method struct {
	Name    string
	Index   int
	Inputs  []Param
	Outputs []Param
}

type Param struct {
	Type  reflect.Type
	Name  string
	Index int
}
