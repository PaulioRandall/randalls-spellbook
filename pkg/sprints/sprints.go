package sprints

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	OptionShowUnexported = "OptionShowUnexported"
)

var (
	ErrBadIndentCount = fmt.Errorf(
		"Integer must follow OptionIndent",
	)
)

type fmtCtx struct {
	unexported bool
}

func parseOptions(options []string) fmtCtx {
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

type strb struct {
	builder strings.Builder
}

func (sb *strb) fmt(msg string, args ...any) {
	s := fmt.Sprintf(msg, args...)
	sb.builder.WriteString(s)
}

func (sb *strb) lnfmt(msg string, args ...any) {
	s := fmt.Sprintf(msg, args...)
	sb.builder.WriteRune('\n')
	sb.builder.WriteString(s)
}

func (sb *strb) lnstr(s string) {
	sb.builder.WriteRune('\n')
	sb.builder.WriteString(s)
}

func (sb *strb) String() string {
	return sb.builder.String()
}

// Print stringifies the object and prints it to terminal.
func Print(object any, options ...string) {
	sb := strb{}
	stringifyObject(
		&sb,
		object,
		parseOptions(options),
	)
	fmt.Print(sb.String())
}

// Println stringifies the object and prints it to
// terminal followed by a linefeed.
func Println(object any, options ...string) {
	sb := strb{}
	stringifyObject(
		&sb,
		object,
		parseOptions(options),
	)
	fmt.Println(sb.String())
}

// String formats an objects into a string form similar
// to the object's definition or instantiation. If a
// non-struct kind is passed then panic ensues.
func String(object any, options ...string) string {
	sb := strb{}
	stringifyObject(
		&sb,
		object,
		parseOptions(options),
	)
	return sb.String()
}

func stringifyObject(
	sb *strb,
	object any,
	ctx fmtCtx,
) {
	val := reflect.ValueOf(object)
	checkObjectType(val)

	sb.fmt("%s {", val.Type().Name())
	writeFields(sb, val, ctx)
	sb.lnstr("}")
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
	sb *strb,
	structVal reflect.Value,
	ctx fmtCtx,
) {
	for field, fieldVal := range structVal.Fields() {
		if !field.IsExported() && !ctx.unexported {
			continue
		}

		sb.lnfmt(
			"\t%s: %v,",
			field.Name,
			getPrintableValue(
				field,
				fieldVal,
			),
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
