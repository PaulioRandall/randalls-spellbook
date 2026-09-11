package nidoking

import (
	"fmt"
)

func ExampleTemplate() {
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
		"	{{filter_column}} IN [{{q_marks}}]",
	).
		Join("columns", ",", columns...).
		Fmt("table", "players").
		Fmt("filter_column", columns[0]).
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
