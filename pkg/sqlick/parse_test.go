package sqlick

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parse_1(t *testing.T) {
	// Happy path.
	// Public fields parsed.
	// Private fields ignored.

	type TestEntity struct {
		Alice   string
		Bob     int
		Charlie float64
		dave    string
	}

	v := TestEntity{}

	act, e := Parse(v)
	require.Equal(t, nil, e)

	exp := SqlickTable{
		GoName: "TestEntity",
		Columns: []SqlickColumn{
			SqlickColumn{
				GoName: "Alice",
				GoType: "string",
			},
			SqlickColumn{
				GoName: "Bob",
				GoType: "int",
			},
			SqlickColumn{
				GoName: "Charlie",
				GoType: "float64",
			},
		},
	}

	require.Equal(t, exp, act)
}

func Test_Parse_2(t *testing.T) {
	// Error when object is not a struct.
	var v int
	_, e := Parse(v)
	require.ErrorIs(t, e, ErrNotStruct)
}

func Test_Parse_3(t *testing.T) {
	// Error when column has unsupported Go type.

	type TestEntity struct {
		Alice error
	}

	v := TestEntity{}
	_, e := Parse(v)
	require.ErrorIs(t, e, ErrBadFieldType)
}

func Test_Parse_4(t *testing.T) {
	// Error when struct has no public fields.

	type TestEntity struct {
		private int
	}

	v := TestEntity{}
	_, e := Parse(v)
	require.ErrorIs(t, e, ErrMissingFields)
}
