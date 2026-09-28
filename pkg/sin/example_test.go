package sin

import (
	"errors"
	"fmt"
)

func Example() {
	ErrNotFound := Template("%s not found")
	ErrAccessDatabase := Err("Database access error")

	findPlayer := func(playerName string) error {
		if playerName == "" {
			return Err("Empty player name passed").
				WrapIn(ErrNotFound).
				Fmt("Player")
		}

		// ...

		err := errors.New("SQL specific error")
		if err != nil {
			return Fmt(
				"Error while searching for player '%s'",
				playerName,
			).
				Wrap(err).
				WrapIn(ErrAccessDatabase)
		}

		// ...

		return nil
	}

	err := findPlayer("")
	fmt.Println(AsStack(err, false))

	err = findPlayer("Fred")
	fmt.Println(AsStack(err, false))
	// Output:
	// Player not found
	//	⤷ Empty player name passed
	// Database access error
	//	⤷ Error while searching for player 'Fred'
	//	⤷ SQL specific error
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

	stack := AsStack(causeWrappedInSymptom, false)
	fmt.Println(stack)
	// Output:
	// The symptom
	//	⤷ The cause
}

func ExampleCurse_WrapIn_template() {
	ErrTable := Template("In table '%s'")
	ErrDatabase := Template("In database '%s'")

	columnError := Err("Bad column").
		WrapIn(ErrTable).Fmt("characters").
		WrapIn(ErrDatabase).Fmt("players")

	fmt.Println(columnError.AsStack(true))
	// Output:
	// Bad column
	//	⤤ In table 'characters'
	//	⤤ In database 'players'
}

func ExampleTemplateCurse() {
	ErrTable := Template("In table '%s'")
	err := errors.New("Bad column")

	columnError := ErrTable.Fmt("characters").Wrap(err)

	fmt.Println(columnError.AsStack(false))
	// Output:
	// In table 'characters'
	//	⤷ Bad column
}

func ExampleAsStack() {
	err1 := fmt.Errorf("Root cause")
	err2 := fmt.Errorf("Cause: %w", err1)
	err3 := fmt.Errorf("Top level: %w", err2)

	errStack := AsStack(err3, false)
	fmt.Println(errStack)

	errStack = AsStack(err3, true)
	fmt.Println(errStack)
	// Output:
	// Top level
	//	⤷ Cause
	//	⤷ Root cause
	// Root cause
	//	⤤ Cause
	//	⤤ Top level
}

func ExampleAsRawStack() {
	err1 := fmt.Errorf("Root cause")
	err2 := fmt.Errorf("Cause: %w", err1)
	err3 := fmt.Errorf("Top level: %w", err2)

	errStack := AsRawStack(err3, false)
	fmt.Println(errStack)

	errStack = AsRawStack(err3, true)
	fmt.Println(errStack)
	// Output:
	// Top level: Cause: Root cause
	//	⤷ Cause: Root cause
	//	⤷ Root cause
	// Root cause
	//	⤤ Cause: Root cause
	//	⤤ Top level: Cause: Root cause
}

func ExampleStackError() {
	err1 := fmt.Errorf("Root cause")
	err2 := fmt.Errorf("Cause: %w", err1)
	err3 := fmt.Errorf("Top level: %w", err2)

	st := Stack(err3, true, true)

	fmt.Println(st.Error())
	// Output:
	// Root cause
	//	⤤ Cause: Root cause
	//	⤤ Top level: Cause: Root cause
}
