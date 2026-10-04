package wizzard

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
	require.ErrorIs(t, e, ErrForTable)
	require.ErrorIs(t, e, ErrForModel)
	require.ErrorIs(t, e, ErrNotStruct)
}

func Test_Parse_2(t *testing.T) {
	// When the parsing an object into a model
	// if the model has no fields
	// then the model with have no fields

	type TestDummy struct{}

	model, e := Parse(TestDummy{})
	require.NoError(t, e)

	exp := Model{
		GoType:  ref.TypeOf(TestDummy{}),
		GoName:  "TestDummy",
		SqlName: "TestDummy",
		Props:   nil,
	}

	require.Equal(t, exp, model)
}

func Test_Parse_3(t *testing.T) {
	// When the parsing an object into a model
	// if object is actually a pointer or pointer to a
	// pointer to the object
	// then the object is dereferenced into the concrete
	// type and that is parsed as the model instead

	ptr := &Dummy{}
	ptrPtr := &ptr
	model, e := Parse(ptrPtr)
	require.NoError(t, e)

	exp := modelOfDummy()

	require.Equal(t, exp, model)
}

func Test_Parse_4(t *testing.T) {
	// When the parsing an object into a model
	// if the model only contains unexported fields
	// then the model with have no fields

	type TestDummy struct {
		unexported int
	}

	model, e := Parse(TestDummy{})
	require.NoError(t, e)

	require.Equal(t, 0, len(model.Props))
}

func Test_Parse_5(t *testing.T) {
	// When the parsing an object into a model
	// if an exported field has an unsupported Go kind
	// then a named error is returned

	type TestDummy struct {
		Id *int
	}

	_, e := Parse(TestDummy{})
	require.ErrorIs(t, e, ErrForTable)
	require.ErrorIs(t, e, ErrForModel)
	require.ErrorIs(t, e, ErrForField)
	require.ErrorIs(t, e, ErrUnsupportedType)
}

func Test_Parse_6(t *testing.T) {
	// When the parsing an object into a model
	// the first exported field will be flagged as the
	// primary key

	type TestDummy struct {
		ignored bool
		Id      int
	}

	model, e := Parse(TestDummy{})
	require.NoError(t, e)

	exp := modelOfDummy().Props[0]
	exp.GoIndex = 1

	require.Equal(t, exp, model.Props[0])
}

func Test_Parse_7(t *testing.T) {
	// When the parsing an object into a model
	// if the object has multiple exported fields
	// then all those fields are parsed into columns
	// within the model

	type TestDummy struct {
		ignored1 bool
		Id       int
		ignored2 []bool
		Name     string
		ignored3 *bool
		Rating   float64
	}

	model, e := Parse(TestDummy{})
	require.NoError(t, e)

	exp := modelOfDummy().Props
	exp[0].GoIndex = 1
	exp[1].GoIndex = 3
	exp[2].GoIndex = 5

	require.Equal(t, exp, model.Props)
}

func Test_ParseAs_1(t *testing.T) {
	// When the parsing an object into a model
	// then the model's SqlName is the name passed
	// and the model's GoName is the name of the Go type

	model, e := ParseAs("AliasName", Dummy{})
	require.NoError(t, e)

	exp := modelOfDummy()
	exp.SqlName = "AliasName"

	require.Equal(t, exp, model)
}

func Test_Table_KeyProp_1(t *testing.T) {
	// When getting the primary key column
	// if the model has multiple fields
	// the primary key column is found and returned

	model, e := Parse(Dummy{})
	require.NoError(t, e)
	require.Equal(t, model.Props[0], model.KeyProp())
}

func Test_Model_Create_1(t *testing.T) {
	// When the creating a table
	// then the table is created

	db := createTestDb(t)
	defer db.Close()

	model := modelOfDummy()
	model.Create(db)

	requireDummyTableExists(t, db)
}

func Test_Model_Upsert_1(t *testing.T) {
	// When upserting data
	// if the data row doesn't exist
	// then the data is inserted
	// when upserting again
	// then the data is updated

	db := createDummyTestDb(t)
	defer db.Close()

	model := modelOfDummy()

	data := Dummy{
		Id:     1,
		Name:   "Alice",
		Rating: 1.11,
	}
	e := model.Upsert(db, data)
	require.NoError(t, e)

	requireDummyTableContains(t, db, data)

	newData := Dummy{
		Id:     1,
		Name:   "Bob",
		Rating: 2.22,
	}
	e = model.Upsert(db, newData)
	require.NoError(t, e)

	requireDummyTableContains(t, db, newData)
}

func Test_Model_Insert_Update_1(t *testing.T) {
	// When inserting data
	// then the data is inserted
	// when updating that data
	// then the data is updated

	db := createDummyTestDb(t)
	defer db.Close()

	model := modelOfDummy()

	data := Dummy{
		Id:     1,
		Name:   "Alice",
		Rating: 1.11,
	}
	e := model.Insert(db, data)
	require.NoError(t, e)

	requireDummyTableContains(t, db, data)

	newData := Dummy{
		Id:     1,
		Name:   "Bob",
		Rating: 2.22,
	}
	e = model.Update(db, newData)
	require.NoError(t, e)

	requireDummyTableContains(t, db, newData)
}

func Test_Model_Delete_1(t *testing.T) {
	// When deleting data
	// if the data row exists
	// then the data is deleted

	db := createDummyTestDb(t)
	defer db.Close()

	model := modelOfDummy()

	data := Dummy{
		Id:     1,
		Name:   "Alice",
		Rating: 1.11,
	}
	e := model.Upsert(db, data)
	require.NoError(t, e)

	dbData := requireDummyTableContains(t, db)
	require.Equal(t, 1, len(dbData))

	e = model.Delete(db, "Id = ?", 1)
	require.NoError(t, e)

	dbData = requireDummyTableContains(t, db)
	require.Equal(t, 0, len(dbData))
}

func Test_Model_Select_1(t *testing.T) {
	// When selcting data
	// if the where clause is empty
	// then all rows are returned

	db := createDummyTestDb(t)
	defer db.Close()

	model := modelOfDummy()

	data := []Dummy{
		Dummy{
			Id:     1,
			Name:   "Alice",
			Rating: 1.11,
		},
		Dummy{
			Id:     2,
			Name:   "Bob",
			Rating: 2.22,
		},
		Dummy{
			Id:     3,
			Name:   "Charlie",
			Rating: 3.33,
		},
	}

	for _, d := range data {
		e := model.Upsert(db, d)
		require.NoError(t, e)
	}

	results, e := model.Select[Dummy](db, "")
	require.NoError(t, e)
	require.Equal(t, data, results)
}

func Test_Model_Select_2(t *testing.T) {
	// When selecting data
	// if the where clause has a statement filtering results
	// and the where statement contains a parameter
	// and an argument is provided
	// then a filtered set of results are returned

	db := createDummyTestDb(t)
	defer db.Close()

	model := modelOfDummy()

	data := []Dummy{
		Dummy{
			Id:     1,
			Name:   "Alice",
			Rating: 1.11,
		},
		Dummy{
			Id:     2,
			Name:   "Bob",
			Rating: 2.22,
		},
		Dummy{
			Id:     3,
			Name:   "Charlie",
			Rating: 3.33,
		},
	}

	for _, d := range data {
		e := model.Upsert(db, d)
		require.NoError(t, e)
	}

	results, e := model.Select[Dummy](db, "Rating > ?", 2)
	require.NoError(t, e)

	exp := data[1:] // Bob and Charlie only.
	require.Equal(t, exp, results)
}

func Test_Model_SelectFirst_1(t *testing.T) {
	// When selecting a specific row
	// if the where clause is empty
	// then the first row of the table is returned

	db := createDummyTestDb(t)
	defer db.Close()

	model := modelOfDummy()

	data := []Dummy{
		Dummy{
			Id:     1,
			Name:   "Alice",
			Rating: 1.11,
		},
		Dummy{
			Id:     2,
			Name:   "Bob",
			Rating: 2.22,
		},
		Dummy{
			Id:     3,
			Name:   "Charlie",
			Rating: 3.33,
		},
	}

	for _, d := range data {
		e := model.Upsert(db, d)
		require.NoError(t, e)
	}

	result, found, e := model.SelectFirst[Dummy](db, "")
	require.NoError(t, e)
	require.Equal(t, true, found)
	require.Equal(t, data[0], result)
}

func Test_Model_SelectFirst_2(t *testing.T) {
	// When selecting a specific row
	// if the where clause has a statement filtering results
	// and the where statement contains a parameter
	// and an argument is provided
	// then the first item of filtered result set is returned

	db := createDummyTestDb(t)
	defer db.Close()

	model := modelOfDummy()

	data := []Dummy{
		Dummy{
			Id:     1,
			Name:   "Alice",
			Rating: 1.11,
		},
		Dummy{
			Id:     2,
			Name:   "Bob",
			Rating: 2.22,
		},
		Dummy{
			Id:     3,
			Name:   "Charlie",
			Rating: 3.33,
		},
	}

	for _, d := range data {
		e := model.Upsert(db, d)
		require.NoError(t, e)
	}

	result, found, e := model.SelectFirst[Dummy](
		db,
		"Id = ?",
		2,
	)
	require.NoError(t, e)
	require.Equal(t, true, found)
	require.Equal(t, data[1], result)
}

func Test_Model_SelectFirst_3(t *testing.T) {
	// When selecting a specific row
	// if the where clause has a statement filtering results
	// and the where statement contains a parameter
	// and an argument is provided
	// but where clause filters all results
	// then the found flag is returned false

	db := createDummyTestDb(t)
	defer db.Close()

	model := modelOfDummy()

	data := []Dummy{
		Dummy{
			Id:     1,
			Name:   "Alice",
			Rating: 1.11,
		},
		Dummy{
			Id:     2,
			Name:   "Bob",
			Rating: 2.22,
		},
		Dummy{
			Id:     3,
			Name:   "Charlie",
			Rating: 3.33,
		},
	}

	for _, d := range data {
		e := model.Upsert(db, d)
		require.NoError(t, e)
	}

	_, found, e := model.SelectFirst[Dummy](
		db,
		"Id = ?",
		4,
	)
	require.NoError(t, e)
	require.Equal(t, false, found)
}

func Test_Model_Drop_1(t *testing.T) {
	// When dropping a table
	// then the table is removed

	db := createDummyTestDb(t)
	defer db.Close()

	requireDummyTableExists(t, db)

	modelOfDummy().Drop(db)

	requireDummyTableDoesNotExist(t, db)
}
