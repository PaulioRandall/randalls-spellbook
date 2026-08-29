package datastore

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestEntity struct {
	Alice   string
	Bob     int
	Charlie float64
	dave    string
}

func Test_parseDbTable_1(t *testing.T) {
	v := TestEntity{}

	act, e := parseDbTable(v)
	require.Equal(t, nil, e)

	exp := DbTable{
		Name: "TestEntity",
		Columns: []DbColumn{
			DbColumn{
				Name:     "Alice",
				DataType: "string",
			},
			DbColumn{
				Name:     "Bob",
				DataType: "int",
			},
			DbColumn{
				Name:     "Charlie",
				DataType: "float64",
			},
		},
	}

	require.Equal(t, exp, act)
}
