package nidorino

/*
func Example() {
	columns := []string{
		"name",
		"level",
		"role",
	}

	s := Given(`
		SELECT
			{{columns}}
		FROM
			{{table}}
		WHERE
			{{id_column}} = ?
	`).
		Join("columns", ",", columns...).
		Fmt("table", "players").
		Fmt("id_column", columns[0]).
		String()

	exp := `
		SELECT
			name,
			level,
			role
		FROM
			players
		WHERE
			name = ?
	`

	Println(s)
	// Output:
	//	SELECT
	//		name,
	//		level,
	//		role
	//	FROM
	//		players
	//	WHERE
	//		name = ?
}
*/
