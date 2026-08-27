package spellbook

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func mockIncantation(data any, err error) Incantation {
	return func(_ any, _ Effect) Effect {
		return Choose(data, err)
	}
}

func requireSubset[T any](
	t *testing.T,
	list, subset []T,
) {
	for i, _ := range subset {
		require.Equal(
			t,
			reflect.ValueOf(list[i]).Pointer(),
			reflect.ValueOf(subset[i]).Pointer(),
		)
	}
}

func Test_Enscribe(t *testing.T) {
	spellbook := Conjure()
	incant := mockIncantation(nil, nil)

	spellbook.Enscribe("spellname", incant)

	exp := []Incantation{
		incant,
	}

	requireSubset(t, exp, spellbook.spells["spellname"])
}

func Test_Describe(t *testing.T) {
	spellbook := Conjure()
	incantA := mockIncantation("A", nil)
	incantB := mockIncantation("B", nil)

	spellbook.Transcribe(
		"spellname",
		incantA,
		incantB,
	)

	act := spellbook.Describe("spellname")
	exp := Spell{
		incantA,
		incantB,
	}

	requireSubset(t, exp, act)
}

func Test_Cast_1(t *testing.T) {
	spellbook := Conjure()
	incant := mockIncantation("data", nil)

	spellbook.Enscribe("spellname", incant)

	act := spellbook.Cast("spellname", nil, "info")
	expA := &effect{result: "info", prior: nil}
	expB := &effect{result: "data", prior: expA}

	require.Equal(t, expB, act)
}

func Test_Cast_2(t *testing.T) {
	spellbook := Conjure()
	incant := mockIncantation("data", nil)

	spellbook.Transcribe(
		"spellname",
		incant,
		incant,
		incant,
	)

	act := spellbook.Cast("spellname", nil, nil)
	expA := &effect{result: nil, prior: nil}
	expB := &effect{result: "data", prior: expA}
	expC := &effect{result: "data", prior: expB}
	expD := &effect{result: "data", prior: expC}

	require.Equal(t, expD, act)
}

func Test_Cast_3(t *testing.T) {
	spellbook := Conjure()

	incant := mockIncantation("data", nil)
	spellbook.Enscribe("spellname", incant)

	var act Effect
	exp := &effect{
		result: "data",
		prior:  &effect{},
	}

	act = spellbook.Cast("spellname", nil, nil)
	require.Equal(t, exp, act)

	act = spellbook.Cast("spellname", nil, nil)
	require.Equal(t, exp, act)

	act = spellbook.Cast("spellname", nil, nil)
	require.Equal(t, exp, act)
}

func Test_Cast_4(t *testing.T) {
	spellbook := Conjure()

	err := errors.New("error")
	incant := mockIncantation(nil, err)
	spellbook.Enscribe("spellname", incant)

	act := spellbook.Cast("spellname", nil, "info")
	expA := &effect{result: "info", prior: nil}
	expB := &effect{error: err, prior: expA}

	require.Equal(t, expB, act)
}

func Test_Cast_5(t *testing.T) {
	spellbook := Conjure()
	incantA := mockIncantation("A", nil)
	incantB := mockIncantation("B", nil)
	incantC := mockIncantation("C", nil)

	spellbook.Transcribe(
		"spellname",
		incantA,
		incantB,
		incantC,
	)

	act := spellbook.Cast("spellname", nil, 0)
	exp0 := &effect{result: 0, prior: nil}
	expA := &effect{result: "A", prior: exp0}
	expB := &effect{result: "B", prior: expA}
	expC := &effect{result: "C", prior: expB}

	require.Equal(t, expC, act)
}

func Test_JsonDemystifyer_1(t *testing.T) {
	type testObject struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	var object testObject
	incant := JsonDemystifyer(&object)

	input := &effect{
		result: []byte(`{
		  "name": "Alice",
		  "age": 24
	  }`),
	}

	act := incant(nil, input)
	exp := &effect{
		result: &testObject{
			Name: "Alice",
			Age:  24,
		},
	}

	require.Equal(t, exp, act)
}
