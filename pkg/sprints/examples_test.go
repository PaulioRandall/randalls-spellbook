package sprints

import (
	"fmt"
)

func Example() {
	type Types struct {
		String     string
		Int        int
		Int32      int32
		Int64      int64
		Float32    float32
		Float64    float64
		unexported string
	}

	t := Types{
		String:     "text",
		Int:        69,
		Int32:      69,
		Int64:      69,
		Float32:    69.69,
		Float64:    69.69,
		unexported: "Alright then, keep your secrets",
	}

	s := Sprints(t)
	fmt.Print(s)
	// Output:
	// type Types struct {
	// 	String: "text",
	// 	Int: int(69),
	// 	Int32: int32(69),
	// 	Int64: int64(69),
	// 	Float32: float32(69.69),
	// 	Float64: float64(69.69),
	// 	unexported: <unexported>,
	// }
}
