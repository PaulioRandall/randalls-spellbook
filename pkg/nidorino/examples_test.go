package nidorino

import (
	"fmt"
)

func Example() {
	columns := []string{
		"name",
		"level",
		"role",
	}

	s := Lines(
		"SELECT",
		"	{{columns}}",
		"FROM",
		"	{{table}}",
		"WHERE",
		"	{{id_column}} = ?",
	).
		Join("columns", ",", columns...).
		Fmt("table", "players").
		Fmt("id_column", columns[0]).
		String()

	fmt.Println(s)
	// Output:
	// SELECT
	//	name,
	//	level,
	//	role
	// FROM
	//	players
	// WHERE
	//	name = ?
}
