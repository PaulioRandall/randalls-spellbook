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
	val, prefix := dereference(reflect.ValueOf(object))
	checkObjectType(val)

	sb := strings.Builder{}

	s := fmt.Sprintf("%s%s {", prefix, val.Type().Name())
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

	var result string
	var prefix string

	val, prefix = dereference(val)
	typ := val.Type()

	switch typ.Kind() {
	case reflect.Struct:
		result = fmtContainerName(val)
	case reflect.Array, reflect.Slice:
		result = fmt.Sprintf(
			"[%d]%s",
			val.Len(),
			fmtCollectionName(val, typ.Elem()),
		)
	case reflect.String:
		result = fmtString(val)
	default:
		result = fmt.Sprintf(
			"%s(%v)",
			typ.Name(),
			val.Interface(),
		)
	}

	return prefix + result
}

func dereference(
	val reflect.Value,
) (reflect.Value, string) {
	derefCount := 0

	for val.Type().Kind() == reflect.Ptr {
		derefCount++

		if val.IsNil() {
			// Create zero value so we have a real value.
			val = reflect.New(val.Type().Elem())
		}

		val = val.Elem()
	}

	prefix := strings.Repeat("*", derefCount)
	return val, prefix
}

func fmtString(val reflect.Value) string {
	const maxStringLength = 30
	s, _ := val.Interface().(string)
	s = strings.Replace(s, "\n", "\\n", -1)

	r := []rune(s)
	if len(r) > maxStringLength {
		s = string(r[:maxStringLength]) + "..."
	}

	return `"` + string(s) + `"`
}

func fmtContainerName(val reflect.Value) string {
	if val.IsZero() {
		return val.Type().Name() + "{}"
	}
	return val.Type().Name() + "{...}"
}

func fmtCollectionName(
	val reflect.Value,
	elemTyp reflect.Type,
) string {
	if val.IsZero() || val.IsNil() || val.Len() == 0 {
		return elemTyp.Name() + "{}"
	}
	return elemTyp.Name() + "{...}"
}
