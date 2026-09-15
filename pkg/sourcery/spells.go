package sourcery

import (
	"fmt"
	"reflect"
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
		panic("All spells must return nothing, an error, or a value and an error")
	}

	if len(outputs) > 0 {
		lastTyp := outputs[len(outputs)-1]
		errTyp := reflect.TypeOf((*error)(nil)).Elem()

		if !lastTyp.Implements(errTyp) {
			panic("All spells must return (error) or (T, error)")
		}
	}

	return inputs, outputs
}

// Invoke calls thee spell returning the result. The args
// must match the types and order of Spell.Accepts or an
// error is returned.
func (sp Spell) Invoke(args ...any) (any, error) {
	lenArgs := len(args)
	lenIns := len(sp.Accepts)

	if lenArgs != lenIns {
		return nil, fmt.Errorf(
			"'%s' requires %d arguments, you gave me %d",
			sp.Name,
			lenIns,
			lenArgs,
		)
	}

	params := make([]reflect.Value, lenArgs, lenArgs)

	for i, v := range args {
		argVal := reflect.ValueOf(v)
		argTyp := argVal.Type()

		params[i] = argVal
		inTyp := sp.Accepts[i]

		if argTyp != inTyp {
			return nil, fmt.Errorf(
				"'%s' requires %s as argument %d, you gave me %s",
				sp.Name,
				inTyp.Name(),
				i+1,
				argTyp.Name(),
			)
		}
	}

	results := reflect.ValueOf(sp.Func).Call(params)

	switch len(results) {
	case 0:
		return nil, nil
	case 1:
		// Type check done during parsing.
		err, _ := results[0].Interface().(error)
		return nil, err
	default: // 2 return values
		// Type check done during parsing.
		err, _ := results[1].Interface().(error)
		return results[0].Interface(), err
	}
}
