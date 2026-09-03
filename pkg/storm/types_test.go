package storm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SqlickTable_String_1(t *testing.T) {
	table := Table{
		GoName: "Person",
		Columns: []Column{
			Column{
				GoName:  "Name",
				GoType:  strType,
				GoIndex: 0,
			},
			Column{
				GoName:  "Age",
				GoType:  int64Type,
				GoIndex: 1,
			},
		},
	}

	exp := joinLines(
		"Person",
		"  [0] Name: string",
		"  [1] Age: int64",
	)

	require.Equal(t, exp, table.String())
}
