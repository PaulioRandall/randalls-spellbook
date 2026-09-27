package curse

import (
	"errors"
	"fmt"
)

func Example() {
	printStack := func(e error) {
		st := Stack(e, false)
		fmt.Println(st)
	}

	// Creating simple errors
	curseErr := Err("1. Plain error message")
	printStack(curseErr)

	// Formating errors
	curseErr = Fmt("%d. Formatted %s message", 2, "error")
	printStack(curseErr)

	// Plundering errors
	err := errors.New("3. Plundered error message")
	curseErr = Plunder(err)
	printStack(curseErr)

	// Wrapping errors
	err = errors.New("Original error")
	curseErr = Err("4. Wrapper for error")
	curseErr = curseErr.Wrap(err)
	printStack(curseErr)

	// Wrapping curses (simple)
	curseErr = Err("Original curse")
	curseErr = curseErr.WrapErr("5. Wrapper for curse")
	printStack(curseErr)

	// Wrapping curses (format)
	curseErr = Err("Original curse")
	curseErr = curseErr.WrapFmt("%d. Wrapper for curse (%s)", 6, "formatted")
	printStack(curseErr)

	// Wrapping within a Curse or TemplateCurse returns the
	// does the opposite of Wrap, i.e. the symptom wraps the
	// calling cause and is returned.
	symptom := Err("7. Symptom curse")
	cause := Err("Cause curse")
	curseErr = cause.WrapIn(symptom)
	printStack(curseErr)

	// Declaring a prototype curse.
	protoErr := Template("%d. %s prototype curse")
	printStack(protoErr)

	// Using a prototype curse.
	// A warning message is printed via log package if you
	// call Error() on TemplateCurse, but the unformatted error
	// message will still be returned.
	curseErr = protoErr.Fmt(9, "Formatted")
	printStack(curseErr)

	// Comparing against the prototype.
	isErr := curseErr.Is(protoErr)
	fmt.Printf("10. Was curseErr created from protoErr: %v\n", isErr)

	// Getting the error chain as a nicely formatted string.
	// Can be reveresed.
	err1 := fmt.Errorf("Root cause")
	err2 := fmt.Errorf("Cause: %w", err1)
	err3 := fmt.Errorf("Top level: %w", err2)

	errStack := Stack(err3, false)
	fmt.Println("11. " + errStack)

	errStack = Stack(err3, true)
	fmt.Println("12. " + errStack)

	// Output:
	// 1. Plain error message
	// 2. Formatted error message
	// 3. Plundered error message
	// 4. Wrapper for error
	//	⤷ Original error
	// 5. Wrapper for curse
	//	⤷ Original curse
	// 6. Wrapper for curse (formatted)
	//	⤷ Original curse
	// 7. Symptom curse
	//	⤷ Cause curse
	// %d. %s prototype curse
	// 9. Formatted prototype curse
	// 10. Was curseErr created from protoErr: true
	// 11. Top level
	//	⤷ Cause
	//	⤷ Root cause
	// 12. Root cause
	//	⮤ Cause
	//	⮤ Top level
}
