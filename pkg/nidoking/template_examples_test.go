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
		ListJoin("columns", ",", columns...).
		Inline("table", "players").
		Inline("filter_column", columns[0]).
		InlineRepeat("q_marks", ", ", "?", len(values)).
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
