package modtab

import (
	ref "reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parse_1(t *testing.T) {
	// GIVEN non-struct model
	// THEN panic ensues

	_, e := Parse(123)
	require.ErrorIs(t, e, ErrForModel)
	require.ErrorIs(t, e, ErrNotStruct)
}

func Test_Parse_2(t *testing.T) {
	// GIVEN empty model
	// THEN  ModelTable fields match expected values
	// AND   ModelTable.Columns is empty

	type TestModel struct{}

	table, e := Parse(TestModel{})
	require.NoError(t, e)

	exp := ModelTable{
		GoType:  ref.TypeOf(TestModel{}),
		GoName:  "TestModel",
		SqlName: "TestModel",
		Columns: nil,
	}

	require.Equal(t, exp, table)
}

func Test_Parse_3(t *testing.T) {
	// GIVEN model being pointed to
	// THEN  derefences and returns expect ModelTable

	type TestModel struct{}

	ptr := &TestModel{}
	ptrPtr := &ptr
	table, e := Parse(ptrPtr)
	require.NoError(t, e)

	exp := ModelTable{
		GoType:  ref.TypeOf(TestModel{}),
		GoName:  "TestModel",
		SqlName: "TestModel",
		Columns: nil,
	}

	require.Equal(t, exp, table)
}

func Test_Parse_4(t *testing.T) {
	// GIVEN model with single unexported field
	// THEN  ModelTable.Columns is empty

	type TestModel struct {
		id int
	}

	table, e := Parse(TestModel{})
	require.NoError(t, e)

	require.Equal(t, 0, len(table.Columns))
}

func Test_Parse_5(t *testing.T) {
	// GIVEN model with exported field of unsupported type
	// THEN  panic ensues

	type TestModel struct {
		Id *int
	}

	_, e := Parse(TestModel{})
	require.ErrorIs(t, e, ErrForModel)
	require.ErrorIs(t, e, ErrUnsupportedType)
}

func Test_Parse_6(t *testing.T) {
	// GIVEN model with single exported field
	// THEN  ModelTable.Columns contains ID field

	type TestModel struct {
		ignored bool
		Id      int
	}

	table, e := Parse(TestModel{})
	require.NoError(t, e)

	exp := []ModelColumn{
		ModelColumn{
			GoType:     ref.TypeOf(int(0)),
			GoIndex:    1,
			GoName:     "Id",
			SqlType:    "INTEGER",
			SqlName:    "Id",
			SqlDefault: int(0),
			PrimaryKey: true,
		},
	}

	require.Equal(t, exp, table.Columns)
}

func Test_Parse_7(t *testing.T) {
	// GIVEN model with multiple exported fields
	// THEN  ModelTable.Columns contains all exported fields

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

	exp := []ModelColumn{
		ModelColumn{
			GoType:     ref.TypeOf(int(0)),
			GoIndex:    1,
			GoName:     "Id",
			SqlType:    "INTEGER",
			SqlName:    "Id",
			SqlDefault: int(0),
			PrimaryKey: true,
		},
		ModelColumn{
			GoType:     ref.TypeOf(""),
			GoIndex:    3,
			GoName:     "Name",
			SqlType:    "TEXT",
			SqlName:    "Name",
			SqlDefault: "''",
			PrimaryKey: false,
		},
		ModelColumn{
			GoType:     ref.TypeOf(float64(0)),
			GoIndex:    5,
			GoName:     "Value",
			SqlType:    "REAL",
			SqlName:    "Value",
			SqlDefault: float64(0),
			PrimaryKey: false,
		},
	}

	require.Equal(t, exp, table.Columns)
}

func Test_Table_PkCol_1(t *testing.T) {
	// GIVEN model with multiple exported fields
	// THEN  returns expected primary key column

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

	pkCol := table.PkCol()

	exp := ModelColumn{
		GoType:     ref.TypeOf(int(0)),
		GoIndex:    1,
		GoName:     "Id",
		SqlType:    "INTEGER",
		SqlName:    "Id",
		SqlDefault: int(0),
		PrimaryKey: true,
	}

	require.Equal(t, exp, pkCol)
	require.Equal(t, table.Columns[0], pkCol)
}
