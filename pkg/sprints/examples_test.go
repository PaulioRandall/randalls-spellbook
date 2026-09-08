package sprints

import (
	"fmt"
)

func Example() {
	// Nested and embedded objects are not parsed. Only the
	// struct name is presented.
	//
	// Think about it, the result would be extremely hard
	// to read given lots of embedding, nesting, or deep
	// nesting. No, String or Print the object separately.
	type Nested struct {
		Meh string
	}

	type Types struct {
		String       string
		Int          int
		Uint8        uint8
		Float64      float64
		Nested       Nested
		ArrayOrSlice []string
		Ptr          *bool
		PtrPtr       **bool
		unexported   string
	}

	nonPtr := true
	ptr := &nonPtr
	ptrPtr := &ptr

	object := Types{
		String:  "text",
		Int:     69,
		Uint8:   69,
		Float64: 69.69,
		Nested: Nested{
			Meh: "Blah",
		},
		ArrayOrSlice: []string{
			"One",
			"Two",
			"Three",
		},
		Ptr:        ptr,
		PtrPtr:     ptrPtr,
		unexported: "Alright then, keep your secrets",
	}

	s := String(object)

	// Alternatively use Println(object)
	fmt.Println(s)
	// Output:
	// Types {
	// 	String: "text",
	// 	Int: int(69),
	// 	Uint8: uint8(69),
	// 	Float64: float64(69.69),
	//	Nested: Nested{...},
	//	ArrayOrSlice: [3]string{...},
	//	Ptr: *bool(true),
	//	PtrPtr: **bool(true),
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
