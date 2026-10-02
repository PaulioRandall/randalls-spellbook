package modtab

import (
	ref "reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parse_1(t *testing.T) {
	// When the parsing an object into a model
	// if the model is not a Go struct kind (reflect.Struct)
	// then a named error is returned

	_, e := Parse(123)
	require.ErrorIs(t, e, ErrForModel)
	require.ErrorIs(t, e, ErrNotStruct)
}

func Test_Parse_2(t *testing.T) {
	// When the parsing an object into a model
	// if the model has no fields
	// then the model with have no fields

	type TestModel struct{}

	table, e := Parse(TestModel{})
	require.NoError(t, e)

	exp := Model{
		GoType:  ref.TypeOf(TestModel{}),
		GoName:  "TestModel",
		SqlName: "TestModel",
		Props:   nil,
	}

	require.Equal(t, exp, table)
}

func Test_Parse_3(t *testing.T) {
	// When the parsing an object into a model
	// if object is actually a pointer or pointer to a
	// pointer to the object
	// then the object is dereferenced into the concrete
	// type and that is parsed as the model instead

	type TestModel struct{}

	ptr := &TestModel{}
	ptrPtr := &ptr
	table, e := Parse(ptrPtr)
	require.NoError(t, e)

	exp := Model{
		GoType:  ref.TypeOf(TestModel{}),
		GoName:  "TestModel",
		SqlName: "TestModel",
		Props:   nil,
	}

	require.Equal(t, exp, table)
}

func Test_Parse_4(t *testing.T) {
	// When the parsing an object into a model
	// if the model only contains unexported fields
	// then the model with have no fields

	type TestModel struct {
		unexported int
	}

	table, e := Parse(TestModel{})
	require.NoError(t, e)

	require.Equal(t, 0, len(table.Props))
}

func Test_Parse_5(t *testing.T) {
	// When the parsing an object into a model
	// if an exported field has an unsupported Go kind
	// then a named error is returned

	type TestModel struct {
		Id *int
	}

	_, e := Parse(TestModel{})
	require.ErrorIs(t, e, ErrForModel)
	require.ErrorIs(t, e, ErrForField)
	require.ErrorIs(t, e, ErrUnsupportedType)
}

func Test_Parse_6(t *testing.T) {
	// When the parsing an object into a model
	// the first exported field will be flagged as the
	// primary key

	type TestModel struct {
		ignored bool
		Id      int
	}

	table, e := Parse(TestModel{})
	require.NoError(t, e)

	exp := []Property{
		Property{
			GoType:     ref.TypeOf(int(0)),
			GoIndex:    1,
			GoName:     "Id",
			SqlType:    "INTEGER",
			SqlName:    "Id",
			SqlDefault: int(0),
			IsKey:      true,
		},
	}

	require.Equal(t, exp, table.Props)
}

func Test_Parse_7(t *testing.T) {
	// When the parsing an object into a model
	// if the object has multiple exported fields
	// then all those fields are parsed into columns
	// within the model

	type TestModel struct {
		ignored1 bool
		Id       int
		ignored2 []bool
		Name     string
		ignored3 *bool
		Value    float64
	}

	table, e := Parse(TestModel{})
	require.NoError(t, e)

	exp := []Property{
		Property{
			GoType:     ref.TypeOf(int(0)),
			GoIndex:    1,
			GoName:     "Id",
			SqlType:    "INTEGER",
			SqlName:    "Id",
			SqlDefault: int(0),
			IsKey:      true,
		},
		Property{
			GoType:     ref.TypeOf(""),
			GoIndex:    3,
			GoName:     "Name",
			SqlType:    "TEXT",
			SqlName:    "Name",
			SqlDefault: "''",
			IsKey:      false,
		},
		Property{
			GoType:     ref.TypeOf(float64(0)),
			GoIndex:    5,
			GoName:     "Value",
			SqlType:    "REAL",
			SqlName:    "Value",
			SqlDefault: float64(0),
			IsKey:      false,
		},
	}

	require.Equal(t, exp, table.Props)
}

func Test_Table_PkCol_1(t *testing.T) {
	// When getting the primary key column
	// if the model has multiple fields
	// the primary key column is found and returned

	type TestModel struct {
		ignored1 bool
		Id       int
		ignored2 []bool
		Name     string
		ignored3 *bool
		Value    float64
	}

	table, e := Parse(TestModel{})
	require.NoError(t, e)

	pkCol := table.KeyProp()

	exp := Property{
		GoType:     ref.TypeOf(int(0)),
		GoIndex:    1,
		GoName:     "Id",
		SqlType:    "INTEGER",
		SqlName:    "Id",
		SqlDefault: int(0),
		IsKey:      true,
	}

	require.Equal(t, exp, pkCol)
	require.Equal(t, table.Props[0], pkCol)
}
