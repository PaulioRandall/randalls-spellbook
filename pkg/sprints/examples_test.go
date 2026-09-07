package sprints

import (
	"fmt"
)

func Example() {
	type Nested struct {
		Meh string
	}

	type Types struct {
		String     string
		Int        int
		Int32      int32
		Int64      int64
		Float32    float32
		Float64    float64
		Nested     Nested
		unexported string
	}

	object := Types{
		String:  "text",
		Int:     69,
		Int32:   69,
		Int64:   69,
		Float32: 69.69,
		Float64: 69.69,
		Nested: Nested{
			Meh: "Blah",
		},
		unexported: "Alright then, keep your secrets",
	}

	s := String(object)

	// Alternatively use Println(object)
	fmt.Println(s)
	// Output:
	// Types {
	// 	String: "text",
	// 	Int: int(69),
	// 	Int32: int32(69),
	// 	Int64: int64(69),
	// 	Float32: float32(69.69),
	// 	Float64: float64(69.69),
	//	Nested: Nested{...},
	// }
}

func Example_options() {
	type Types struct {
		Name       string
		Level      int
		unexported string
	}

	object := Types{
		Name:       "Bob",
		Level:      69,
		unexported: "Alright then, keep your secrets",
	}

	Println(
		object,
		OptionShowUnexported,
	)
	// Output:
	// Types {
	// 	Name: "Bob",
	//	Level: int(69),
	// 	unexported: <unexported>,
	// }
}
