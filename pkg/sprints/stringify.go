package sprints

import (
	"fmt"
	"reflect"
	"strings"
)

type fmtCtx struct {
	unexported bool
}

func newFmtCtx(options []string) fmtCtx {
	ctx := fmtCtx{}

	for i := 0; i < len(options); i++ {
		op := options[i]

		if op == OptionShowUnexported {
			ctx.unexported = true
			continue
		}
	}

	return ctx
}

func stringifyObject(
	object any,
	ctx fmtCtx,
) string {
	val := reflect.ValueOf(object)
	checkObjectType(val)

	sb := strings.Builder{}

	s := fmt.Sprintf("%s {", val.Type().Name())
	sb.WriteString(s)
	writeFields(&sb, val, ctx)
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

func writeFields(
	sb *strings.Builder,
	structVal reflect.Value,
	ctx fmtCtx,
) {
	for field, fieldVal := range structVal.Fields() {
		if !field.IsExported() && !ctx.unexported {
			continue
		}

		s := fmt.Sprintf(
			"\t%s: %v,",
			field.Name,
			getPrintableValue(
				field,
				fieldVal,
			),
		)
		sb.WriteRune('\n')
		sb.WriteString(s)
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
	case reflect.Struct:
		return field.Type.Name() + "{...}"
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
