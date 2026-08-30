package sqlick

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SqlickTable_String_1(t *testing.T) {
	table := SqlickTable{
		GoName: "Person",
		Columns: []SqlickColumn{
			SqlickColumn{
				GoName: "Name",
				GoType: "string",
			},
			SqlickColumn{
				GoName: "Age",
				GoType: "int",
			},
		},
	}

	exp := joinLines(
		"Person",
		"  Name: string",
		"  Age: int",
	)

	require.Equal(t, exp, table.String())
}
