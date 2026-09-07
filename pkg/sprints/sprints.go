package sprints

import (
	"fmt"
	"reflect"
	"strings"
)

type strb struct {
	strings.Builder
}

func (sb *strb) fmt(msg string, args ...any) {
	s := fmt.Sprintf(msg, args...)
	sb.WriteString(s)
}

func Sprints(object any) string {
	val := reflect.ValueOf(object)
	checkObjectType(val)

	sb := &strb{}

	sb.fmt("type %s struct {", val.Type().Name())
	writeFields(sb, val)
	sb.WriteString("\n}")

	return sb.String()
}

func checkObjectType(val reflect.Value) {
	k := val.Type().Kind()
	if k != reflect.Struct {
		e := fmt.Errorf(
			"Sprints only stringifies Structs, got '%s'",
			k.String(),
		)
		panic(e)
	}
}

func writeFields(sb *strb, structVal reflect.Value) {
	for field, fieldVal := range structVal.Fields() {
		sb.fmt(
			"\n\t%s: %v,",
			field.Name,
			getPrintableValue(field, fieldVal),
		)
	}
}

func getPrintableValue(
	field reflect.StructField,
	val reflect.Value,
) any {
	if !field.IsExported() {
		return "<unexported>"
	}

	switch field.Type.Kind() {
	case reflect.String:
		// TODO: If string is longer than say 40 bytes
		//       then cut it off and append ...
		return fmt.Sprintf(`"%v"`, val.Interface())
	default:
		return fmt.Sprintf(
			"%s(%v)",
			field.Type.Name(),
			val.Interface(),
		)
	}
}
