package wizzard

import (
	"reflect"

	"github.com/PaulioRandall/randalls-spellbook/pkg/curse"
)

type Wizzard struct {
	cache map[reflect.Type]Struct
}

var global = Wizzard{}

var (
	ErrNotStruct = curse.Proto(
		"Expected kind Struct but got %s",
	)

	ErrBadFieldKind = curse.Proto(
		"%s.%s has invalid field kind %s",
	)
)

func Parse(object any) (Struct, error) {
	return global.Parse(object)
}

func New() *Wizzard {
	return &Wizzard{
		cache: map[reflect.Type]Struct{},
	}
}

func (w *Wizzard) Parse(object any) (Struct, error) {
	structTyp := reflect.TypeOf(object)

	// Check cache first.
	if st, ok := w.cache[structTyp]; ok {
		return st, nil
	}

	// If not a struct then bugger off.
	if structTyp.Kind() != reflect.Struct {
		return Struct{}, ErrNotStruct.Fmt(
			structTyp.Kind().String(),
		)
	}

	fields, e := parseFields(structTyp)
	if e != nil {
		return Struct{}, e
	}

	methods := parseMethods(structTyp)

	st := Struct{
		Type:    structTyp,
		Name:    structTyp.Name(),
		Fields:  fields,
		Methods: methods,
	}

	return st, nil
}

func parseFields(structTyp reflect.Type) ([]Field, error) {
	var fields []Field

	for i := 0; i < structTyp.NumField(); i++ {
		structField := structTyp.Field(i)

		if !structField.IsExported() {
			continue
		}

		fieldTyp := structField.Type
		if !isValidFieldKind(fieldTyp.Kind()) {
			return nil, ErrBadFieldKind.Fmt(
				structTyp.Name(),
				structField.Name,
				fieldTyp.Kind(),
			)
		}

		field := Field{
			Type:  fieldTyp,
			Name:  structField.Name,
			Index: i,
		}

		fields = append(fields, field)
	}

	return fields, nil
}

func isValidFieldKind(needle reflect.Kind) bool {
	haystack := []reflect.Kind{
		reflect.Bool,
		reflect.Int,
		reflect.Int32,
		reflect.Int64,
		reflect.String,
	}

	for _, k := range haystack {
		if needle == k {
			return true
		}
	}

	return false
}

func parseMethods(structTyp reflect.Type) []Method {
	var methods []Method

	for i := 0; i < structTyp.NumMethod(); i++ {
		refMethod := structTyp.Method(i)

		if !refMethod.IsExported() {
			continue
		}

		methodTyp := refMethod.Type
		inputs := parseMethodInputs(methodTyp)
		outputs := parseMethodOutputs(methodTyp)

		method := Method{
			Name:    refMethod.Name,
			Index:   refMethod.Index,
			Inputs:  inputs,
			Outputs: outputs,
		}

		methods = append(methods, method)
	}

	return methods
}

func parseMethodInputs(methodTyp reflect.Type) []Param {
	var params []Param

	// Starts at 1 because first param is 'self' or 'this'.
	for i := 1; i < methodTyp.NumIn(); i++ {
		refIn := methodTyp.In(i)

		p := Param{
			Type:  refIn,
			Name:  refIn.Name(),
			Index: i,
		}

		params = append(params, p)
	}

	return params
}

func parseMethodOutputs(methodTyp reflect.Type) []Param {
	var params []Param

	for i := 0; i < methodTyp.NumOut(); i++ {
		refOut := methodTyp.Out(i)

		p := Param{
			Type:  refOut,
			Name:  refOut.Name(),
			Index: i,
		}

		params = append(params, p)
	}

	return params
}
