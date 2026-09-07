package sprints

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

const (
	OptionShowUnexported = "OptionShowUnexported"
)

type fmtOptions struct {
	unexported bool
}

func parseOptions(options []string) fmtOptions {
	opts := fmtOptions{}

	if slices.Contains(options, OptionShowUnexported) {
		opts.unexported = true
	}

	return opts
}

type strb struct {
	strings.Builder
}

func (sb *strb) fmt(msg string, args ...any) {
	s := fmt.Sprintf(msg, args...)
	sb.WriteString(s)
}

// Print stringifies the object and prints it to terminal.
func Print(object any, options ...string) {
	s := stringifyObject(object, parseOptions(options))
	fmt.Print(s)
}

// Println stringifies the object and prints it to
// terminal followed by a linefeed.
func Println(object any, options ...string) {
	s := stringifyObject(object, parseOptions(options))
	fmt.Println(s)
}

// String formats an objects into a string form similar
// to the object's definition or instantiation. If a
// non-struct kind is passed then panic ensues.
func String(object any, options ...string) string {
	return stringifyObject(
		object,
		parseOptions(options),
	)
}

func stringifyObject(object any, opts fmtOptions) string {
	val := reflect.ValueOf(object)
	checkObjectType(val)

	sb := &strb{}

	sb.fmt("type %s struct {", val.Type().Name())
	writeFields(sb, val, opts)
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
	sb *strb,
	structVal reflect.Value,
	opts fmtOptions,
) {
	for field, fieldVal := range structVal.Fields() {
		if !field.IsExported() && !opts.unexported {
			continue
		}

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
