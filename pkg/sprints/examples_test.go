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
		Empty        []int64
		Nil          []int64
		Ptr          *bool
		PtrPtr       **bool
		unexported   string
	}

	var nonPtr bool = true
	ptr := &nonPtr

	var nilPtr *bool = nil
	ptrNilPtr := &nilPtr

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
		Empty:      []int64{},
		Nil:        nil,
		Ptr:        ptr,
		PtrPtr:     ptrNilPtr,
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
	//	Empty: [0]int64{},
	//	Nil: []int64,
	//	Ptr: *bool(true),
	//	PtrPtr: *⁎bool(false),
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
