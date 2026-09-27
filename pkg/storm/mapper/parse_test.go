package mapper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func typeOf(model any) reflect.Type {
	return reflect.TypeOf(model)
}

func Test_ParseModel_1(t *testing.T) {
	// GIVEN non-struct model
	// THEN panic ensues

	require.Panics(t, func() {
		ParseModel(123)
	})
}

func Test_ParseModel_2(t *testing.T) {
	// GIVEN empty model
	// THEN  Table fields match expected values
	// AND   Table.Columns is empty

	type TestModel struct{}

	table := ParseModel(TestModel{})

	exp := Table{
		GoType:  typeOf(TestModel{}),
		GoName:  "TestModel",
		SqlName: "TestModel",
		Columns: nil,
	}

	require.Equal(t, exp, table)
}

func Test_ParseModel_3(t *testing.T) {
	// GIVEN model being pointed to
	// THEN  derefences and returns expect Table

	type TestModel struct{}

	ptr := &TestModel{}
	ptrPtr := &ptr
	table := ParseModel(ptrPtr)

	exp := Table{
		GoType:  typeOf(TestModel{}),
		GoName:  "TestModel",
		SqlName: "TestModel",
		Columns: nil,
	}

	require.Equal(t, exp, table)
}

func Test_ParseModel_4(t *testing.T) {
	// GIVEN model with single unexported field
	// THEN  Table.Columns is empty

	type TestModel struct {
		id int
	}

	table := ParseModel(TestModel{})
	require.Equal(t, 0, len(table.Columns))
}

func Test_ParseModel_5(t *testing.T) {
	// GIVEN model with exported field of unsupported type
	// THEN  panic ensues

	type TestModel struct {
		Id *int
	}

	require.Panics(t, func() {
		ParseModel(TestModel{})
	})
}

func Test_ParseModel_6(t *testing.T) {
	// GIVEN model with single exported field
	// THEN  Table.Columns contains ID field

	type TestModel struct {
		ignored bool
		Id      int
	}

	table := ParseModel(TestModel{})

	exp := []Column{
		Column{
			GoType:     typeOf(int(123)),
			GoIndex:    1,
			GoName:     "Id",
			SqlType:    "INTEGER",
			SqlName:    "Id",
			PrimaryKey: true,
		},
	}

	require.Equal(t, exp, table.Columns)
}

func Test_ParseModel_7(t *testing.T) {
	// GIVEN model with multiple exported fields
	// THEN  Table.Columns contains all exported fields

	type TestModel struct {
		ignored1 bool
		Id       int
		ignored2 []bool
		Name     string
		ignored3 *bool
		Value    float64
	}

	table := ParseModel(TestModel{})

	exp := []Column{
		Column{
			GoType:     typeOf(int(123)),
			GoIndex:    1,
			GoName:     "Id",
			SqlType:    "INTEGER",
			SqlName:    "Id",
			PrimaryKey: true,
		},
		Column{
			GoType:     typeOf(""),
			GoIndex:    3,
			GoName:     "Name",
			SqlType:    "TEXT",
			SqlName:    "Name",
			PrimaryKey: false,
		},
		Column{
			GoType:     typeOf(float64(1.23)),
			GoIndex:    5,
			GoName:     "Value",
			SqlType:    "REAL",
			SqlName:    "Value",
			PrimaryKey: false,
		},
	}

	require.Equal(t, exp, table.Columns)
}

func Test_Table_PrimaryKeyColumn_1(t *testing.T) {
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

	table := ParseModel(TestModel{})
	pkCol := table.PrimaryKeyColumn()

	exp := Column{
		GoType:     typeOf(int(123)),
		GoIndex:    1,
		GoName:     "Id",
		SqlType:    "INTEGER",
		SqlName:    "Id",
		PrimaryKey: true,
	}

	require.Equal(t, exp, pkCol)
	require.Equal(t, table.Columns[0], pkCol)
}
