package curse

import (
	"errors"
	"fmt"
)

func Example() {
	rootCause := errors.New("Level 1")
	cause := Err("Level 2").Wrap(rootCause)
	cu := Fmt("Level %d", 3).Wrap(cause)

	fmt.Println(cu)
	// Output:
	// Level 3:
	//	Level 2:
	//	Level 1
}

func ExampleNew() {
	cu := New("Error message")

	fmt.Println(cu)
	// Output:
	// Error message
}

func ExampleNewf() {
	cu := Newf("Error message %s", "formatted")

	fmt.Println(cu)
	// Output:
	// Error message formatted
}

func ExampleErr() {
	cu := Err("Error message")

	fmt.Println(cu)
	// Output:
	// Error message
}

func ExampleFmt() {
	cu := Fmt("Error message %s", "formatted")

	fmt.Println(cu)
	// Output:
	// Error message formatted
}

func ExampleCurse_Wrap() {
	cause := errors.New("Root cause")
	cu := New("Error message").Wrap(cause)

	fmt.Println(cu)
	// Output:
	// Error message:
	//	Root cause
}
