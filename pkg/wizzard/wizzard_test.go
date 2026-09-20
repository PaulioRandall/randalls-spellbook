package wizzard

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type Dummy struct {
	Bool       bool
	Int        int
	unexported string
	String     string
}

func (Dummy) NoParams() {
}

func unexportedMethod() {
}

func InputParams(index int, value **string, err error) {
}

func OutputParams() (string, error) {
	return "", nil
}

func nilErrorType() reflect.Type {
	return reflect.TypeOf((*error)(nil)).Elem()
}

func findField(
	object any,
	fieldName string,
) reflect.StructField {
	typ := reflect.TypeOf(object)
	f, ok := typ.FieldByName(fieldName)

	if ok {
		return f
	}

	panic("No such field '" + fieldName + "' in '" + typ.Name() + "'")
}

func findMethod(
	object any,
	methodName string,
) reflect.Method {
	typ := reflect.TypeOf(object)
	m, _ := typ.MethodByName(methodName)
	return m
}

func TestStruct_Parse_1(t *testing.T) {
	// Bad if not a struct.
	w := New()
	_, e := w.Parse(1)
	require.ErrorIs(t, e, ErrNotStruct)
}

func TestStruct_Parse_2(t *testing.T) {
	// Good if Struct has correct Type.
	w := New()

	act, e := w.Parse(Dummy{})
	exp := reflect.TypeOf(Dummy{})

	require.NoError(t, e)
	require.Equal(t, exp, act.Type)
}

func TestStruct_Parse_3(t *testing.T) {
	// Good if Struct has correct Name.
	w := New()

	act, e := w.Parse(Dummy{})
	exp := "Dummy"

	require.NoError(t, e)
	require.Equal(t, exp, act.Name)
}

func TestStruct_Parse_4(t *testing.T) {
	// Bad if field has invalid type kind.
	w := New()

	type Failure struct {
		TestField *int
	}

	_, e := w.Parse(Failure{})
	require.ErrorIs(t, e, ErrBadFieldKind)
}

func TestStruct_Parse_5(t *testing.T) {
	// Good if field has valid type kind.
	w := New()

	type Success struct {
		Bool   bool
		Int    int
		Int32  int32
		Int64  int64
		String string
	}

	_, e := w.Parse(Success{})
	require.NoError(t, e)
}

func TestStruct_Parse_6(t *testing.T) {
	// Good if parses all exported fields.
	w := New()

	act, e := w.Parse(Dummy{})

	expField1 := Field{
		Type:  reflect.TypeOf(bool(false)),
		Name:  "Bool",
		Index: 0,
	}

	expField2 := Field{
		Type:  reflect.TypeOf(int(0)),
		Name:  "Int",
		Index: 1,
	}

	expField3 := Field{
		Type:  reflect.TypeOf(string("")),
		Name:  "String",
		Index: 3,
	}

	require.NoError(t, e)
	require.Equal(t, expField1, act.Fields[0])
	require.Equal(t, expField2, act.Fields[1])
	require.Equal(t, expField3, act.Fields[2])
	require.Equal(t, 3, len(act.Fields))
}

func TestStruct_Parse_7(t *testing.T) {
	// Good if parses exported method with no params.
	w := New()

	act, e := w.Parse(Dummy{})

	exp := Method{
		Name:    "NoParams",
		Index:   0,
		Inputs:  nil,
		Outputs: nil,
	}

	require.NoError(t, e)
	require.Equal(t, exp, act.Methods[0])
}
