package sin

import (
	"errors"
	"fmt"
)

func Example() {
	ErrNotFound := Err("Player not found in database")
	ErrForPlayer := Template("Regarding player '%s'")
	ErrForTable := Template("Regarding table '%s'")

	err := ErrForPlayer.
		Fmt("Bob").
		Wrap(ErrNotFound).
		WrapIn(ErrForTable).
		Fmt("Players")

	fmt.Println(StackString(err))
	// Output:
	// Regarding table 'Players'
	//	⤷ Regarding player 'Bob'
	//	⤷ Player not found in database
}

func ExampleErr() {
	err := Err("Plain and simple Curse")
	fmt.Println(err.Error())
	// Output:
	// Plain and simple Curse
}

func ExampleFmt() {
	err := Fmt("%s %s", "Formatted", "Curse")
	fmt.Println(err.Error())
	// Output:
	// Formatted Curse
}

func ExamplePlunder() {
	err := errors.New("Steal this error message")
	newErr := Plunder(err)
	fmt.Println(newErr.Error())
	// Output:
	// Steal this error message
}

func ExampleCurse_Is() {
	err := errors.New("Standard string error")
	solo := Err("Solo curse")
	template := Template("%s curse")
	derived := template.Fmt("Derived")

	fmt.Printf("solo is err: %v\n", solo.Is(err))
	fmt.Printf("solo is template: %v\n", solo.Is(template))
	fmt.Printf("derived is err: %v\n", derived.Is(err))
	fmt.Printf("derived is solo: %v\n", derived.Is(solo))
	fmt.Printf("derived is template: %v\n", derived.Is(template))
	// Output:
	// solo is err: false
	// solo is template: false
	// derived is err: false
	// derived is solo: false
	// derived is template: true
}

func ExampleCurse_WrapIn() {
	cause := Err("The cause")
	symptom := Err("The symptom")

	causeWrappedInSymptom := cause.WrapIn(symptom)

	stack := StackString(causeWrappedInSymptom)
	fmt.Println(stack)
	// Output:
	// The symptom
	//	⤷ The cause
}

func ExampleCurse_WrapIn_template() {
	ErrTable := Template("In table '%s'")
	ErrDatabase := Template("In database '%s'")

	colError := Err("Bad column").
		WrapIn(ErrTable).Fmt("characters").
		WrapIn(ErrDatabase).Fmt("players")

	fmt.Println(colError.ReverseStackString())
	// Output:
	// Bad column
	//	⮤ In table 'characters'
	//	⮤ In database 'players'
}

func ExampleTemplateCurse() {
	ErrTable := Template("In table '%s'")
	err := errors.New("Bad column")

	colError := ErrTable.Fmt("characters").Wrap(err)

	fmt.Println(colError.StackString())
	// Output:
	// In table 'characters'
	//	⤷ Bad column
}

func ExampleStackString() {
	err1 := Err("Root cause")
	err2 := Err("Cause")
	err3 := Err("Top level")
	cu := err1.WrapIn(err2).WrapIn(err3)

	errStack := StackString(cu)
	fmt.Println(errStack)

	errStack = ReverseStackString(cu)
	fmt.Println(errStack)
	// Output:
	// Top level
	//	⤷ Cause
	//	⤷ Root cause
	// Root cause
	//	⮤ Cause
	//	⮤ Top level
}

func ExampleStackError() {
	err1 := Err("Root cause")
	err2 := Err("Cause")
	err3 := Err("Top level")
	cu := err1.WrapIn(err2).WrapIn(err3)

	st := Stack(cu)

	fmt.Println(st.Error())
	// Output:
	// Top level
	//	⤷ Cause
	//	⤷ Root cause
}
