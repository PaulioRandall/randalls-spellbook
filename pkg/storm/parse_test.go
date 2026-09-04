package storm

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parse_1(t *testing.T) {
	// Happy path.
	// Public fields parsed.
	// Private fields ignored.

	type TestEntity struct {
		Alice   string
		Bob     int64
		dave    string
		Charlie float64
	}

	v := TestEntity{}

	act, e := Parse(v)
	require.Equal(t, nil, e)

	exp := Table{
		GoType: reflect.TypeOf(v),
		GoName: "TestEntity",
		Columns: []Column{
			Column{
				GoName:  "Alice",
				GoType:  strType,
				GoIndex: 0,
			},
			Column{
				GoName:  "Bob",
				GoType:  int64Type,
				GoIndex: 1,
			},
			Column{
				GoName:  "Charlie",
				GoType:  float64Type,
				GoIndex: 3,
			},
		},
	}

	require.Equal(t, exp, act)
}

func Test_Parse_2(t *testing.T) {
	// Error when object is not a struct.
	_, e := Parse(int64Zero)
	require.ErrorIs(t, e, ErrNotStruct)
}

func Test_Parse_3(t *testing.T) {
	// Error when column has unsupported Go type.

	type TestEntity struct {
		Alice error
	}

	v := TestEntity{}
	_, e := Parse(v)
	require.ErrorIs(t, e, ErrBadFieldKind)
}

func Test_Parse_4(t *testing.T) {
	// Error when struct has no public fields.

	type TestEntity struct {
		private int64
	}

	v := TestEntity{}
	_, e := Parse(v)
	require.ErrorIs(t, e, ErrMissFields)
}
