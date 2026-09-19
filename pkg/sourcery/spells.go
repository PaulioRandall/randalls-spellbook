package sourcery

import (
	"reflect"

	"github.com/PaulioRandall/randalls-spellbook/pkg/curse"
)

var (
	ErrInvokeFail = curse.Proto(
		"Failed to invoke %s",
	)

	ErrArgCount = curse.Proto(
		"Require %d arguments, given %d",
	)

	ErrArgTypeMismatch = curse.Proto(
		"Expected %s or compatible type, given %s",
	)
)

// TODO: Clean up and test pkg

type refT = reflect.Type

// Spell represents a function that may be called. The
// function could have any number of arguments but will
// either return 1 or 2 values, and always an error as
// the last return value. Use [NewSpell] to create.
type Spell struct {
	Name    string
	Func    any
	Accepts []refT
	Returns []refT
}

// NewSpell creates a new [Spell]. Spells can have any
// number or arguments but must have 1 or 2 return values.
// Its last return value must always be an error.
func NewSpell(name string, fn any) Spell {
	inputs, outputs := parseFunc(fn)

	return Spell{
		Name:    name,
		Func:    fn,
		Accepts: inputs,
		Returns: outputs,
	}
}

func parseFunc(fn any) ([]refT, []refT) {
	typ := reflect.TypeOf(fn)
	if typ.Kind() != reflect.Func {
		panic(
			"All spells must be functions, you gave me a " +
				typ.Kind().String(),
		)
	}

	inputs := make([]refT, typ.NumIn(), typ.NumIn())
	for i := 0; i < typ.NumIn(); i++ {
		inputs[i] = typ.In(i)
	}

	outputs := make([]refT, typ.NumOut(), typ.NumOut())
	for i := 0; i < typ.NumOut(); i++ {
		outputs[i] = typ.Out(i)
	}

	if len(outputs) > 2 {
		panic("All spells must return either nothing, a value, an error, or a value and an error")
	}

	if len(outputs) == 2 {
		errTyp := reflect.TypeOf((*error)(nil)).Elem()
		if !outputs[1].Implements(errTyp) {
			panic("Second return value my only be an error")
		}
	}

	return inputs, outputs
}

// Invoke calls thee spell returning the result. The args
// must match the types and order of Spell.Accepts or an
// error is returned.
func (sp Spell) Invoke(args ...any) (any, error) {
	e := sp.checkArgsCount(args)
	if e != nil {
		return nil, ErrInvokeFail.Fmt(sp.Name).Wrap(e)
	}

	argVals, e := parseArgs(args, sp.Accepts)
	if e != nil {
		return nil, ErrInvokeFail.Fmt(sp.Name).Wrap(e)
	}

	results := reflect.ValueOf(sp.Func).Call(argVals)
	return handleInvokeResults(results)
}

func (sp Spell) checkArgsCount(args []any) error {
	argsLen := len(args)
	expLen := len(sp.Accepts)

	if argsLen != expLen {
		return ErrArgCount.Fmt(expLen, argsLen)
	}

	return nil
}

func parseArgs(
	args []any,
	expTypes []reflect.Type,
) ([]reflect.Value, error) {
	argVals := make([]reflect.Value, len(args), len(args))

	for i, arg := range args {
		argVal, e := parseArg(arg, expTypes[i])

		if e != nil {
			return nil, curse.Fmt("For argument %d", i).Wrap(e)
		}

		argVals[i] = argVal
	}

	return argVals, nil
}

func parseArg(
	arg any,
	expTyp reflect.Type,
) (reflect.Value, error) {
	argVal := reflect.ValueOf(arg)
	argTyp := argVal.Type()

	if argTyp == expTyp {
		return argVal, nil
	}

	if castVal, ok := castArgValToExpType(argVal, expTyp); ok {
		return castVal, nil
	}

	return argVal, ErrArgTypeMismatch.Fmt(
		expTyp.Name(),
		argTyp.Name(),
	)
}

func castArgValToExpType(
	argVal reflect.Value,
	expTyp reflect.Type,
) (reflect.Value, bool) {
	isFloat64 := argVal.Type().Kind() == reflect.Float64

	if !isFloat64 {
		return argVal, false
	}

	if isFloat64 && expTyp.Kind() == reflect.Int {
		v, _ := argVal.Interface().(float64)
		return reflect.ValueOf(int(v)), true
	}

	if isFloat64 && expTyp.Kind() == reflect.Int64 {
		v, _ := argVal.Interface().(float64)
		return reflect.ValueOf(int64(v)), true
	}

	if isFloat64 && expTyp.Kind() == reflect.Int32 {
		v, _ := argVal.Interface().(float64)
		return reflect.ValueOf(int32(v)), true
	}

	return argVal, false
}

func handleInvokeResults(
	results []reflect.Value,
) (any, error) {
	switch len(results) {
	case 0:
		return nil, nil
	case 1:
		v := results[0].Interface()
		if err, ok := v.(error); ok {
			return nil, err
		}
		return v, nil
	default: // 2 return values
		// Type check done during parsing.
		err, _ := results[1].Interface().(error)
		return results[0].Interface(), err
	}
}
