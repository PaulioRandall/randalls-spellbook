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

	values := []string{
		"Alice",
		"Bob",
		"Charlie",
	}

	s := Lines(
		"SELECT",
		"	{{columns}}",
		"FROM",
		"	{{table}}",
		"WHERE",
		"	{{id_column}} IN [{{q_marks}}]",
	).
		Join("columns", ",", columns...).
		Fmt("table", "players").
		Fmt("id_column", columns[0]).
		FmtRepeat("q_marks", ", ", "?", len(values)).
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
	//	name IN [?, ?, ?]
}
