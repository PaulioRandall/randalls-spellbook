package storm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Table_String_1(t *testing.T) {
	table := Table{
		GoName: "Person",
		Columns: []Column{
			Column{
				GoName:  "Name",
				GoType:  strType,
				GoIndex: 0,
				SqlType: "TEXT",
			},
			Column{
				GoName:  "Age",
				GoType:  int64Type,
				GoIndex: 1,
				SqlType: "INTEGER",
			},
		},
	}

	exp := joinLines(
		"Person",
		"  [0] Name: Go(string) SQL(TEXT)",
		"  [1] Age: Go(int64) SQL(INTEGER)",
	)

	require.Equal(t, exp, table.String())
}
