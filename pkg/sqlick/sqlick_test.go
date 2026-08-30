package sqlick

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SqlickTable_String_1(t *testing.T) {
	table := Table{
		GoName: "Person",
		Columns: []Column{
			Column{
				GoName: "Name",
				GoType: "string",
			},
			Column{
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
