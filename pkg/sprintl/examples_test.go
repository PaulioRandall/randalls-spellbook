package sprintl

import (
	"fmt"
)

func Example() {
	md := Lines(
		"# TODO %s", // Line 1
		"",
		"• %s", // Line 3
	).
		Fmt(1, "Today").
		Rep(3, "Play", "Eat", "Read", "Sleep").
		String()

	fmt.Println(md)
	// Output:
	// # TODO Today
	//
	// • Play
	// • Eat
	// • Read
	// • Sleep
}

// Example usage with lots of parameters to show off
// features. Most query formatting is not this complex.
func Example_complex() {
	props := []any{
		"name",
		"age",
		"job",
	}

	table := "users"

	filters := []any{
		"name = ?",
		"age < ?",
		"job LIKE 'Master %'",
	}

	sql := Lines(
		"SELECT",
		"  %s", // Line 2
		"FROM",
		"  %s", // Line 4
		"WHERE",
		"%s", // Line 6
	).
		Rep(2, props...).
		Join(",").
		Fmt(4, table).
		Rep(6, filters...).
		Marry("  ", "AND ").
		String()

	fmt.Println(sql)
	// Output:
	// SELECT
	//   name,
	//   age,
	//   job
	// FROM
	//   users
	// WHERE
	//   name = ?
	//   AND age < ?
	//   AND job LIKE 'Master %'
}

func ExampleLines() {
	var _ *Sprintl = Lines(
		"SELECT",
		"  %s", // Line 2
		"FROM",
		"  %s", // Line 4
	)
}

func ExampleSprintl_Fmt() {
	sql := Lines(
		"SELECT",
		"  'month_' || '%d' AS %s,", // Line 2
		"  revenue",
		"FROM",
		"  sales_data",
		"WHERE",
		"  year = 2025",
	).
		Fmt(2, 9, "sept").
		String()

	fmt.Println(sql)
	// Output:
	// SELECT
	//   'month_' || '9' AS sept,
	//   revenue
	// FROM
	//   sales_data
	// WHERE
	//   year = 2025
}

func ExampleSprintl_Dup() {
	sql := Lines(
		"SELECT",
		"  name,",
		"  age,",
		"  job",
		"FROM",
		"  users",
		"WHERE",
		"  name IN [",
		"    %s", // Line 9
		"  ]",
	).
		Dup(9, 4, "?").
		Join(",").
		String()

	fmt.Println(sql)
	// Output:
	// SELECT
	//   name,
	//   age,
	//   job
	// FROM
	//   users
	// WHERE
	//   name IN [
	//     ?,
	//     ?,
	//     ?,
	//     ?
	//   ]
}

func ExampleSprintl_Rep() {
	sql := Lines(
		"SELECT",
		"  %s", // Line 2
		"FROM",
		"  users",
		"WHERE",
		"  %s = ?", // Line 6
	).
		Rep(2, "name", "age", "job").
		Join(",").
		Fmt(6, "name").
		String()

	fmt.Println(sql)
	// Output:
	// SELECT
	//   name,
	//   age,
	//   job
	// FROM
	//   users
	// WHERE
	//   name = ?
}

func ExampleSprintl_Rep_args() {
	sql := Lines(
		"SELECT",
		"  %s AS %s", // Line 2
		"FROM",
		"  users",
	).
		Rep(2,
			[]any{"name", "player"},
			[]any{"age", "level"},
			[]any{"job", "role"},
		).
		Join(",").
		String()

	fmt.Println(sql)
	// Output:
	// SELECT
	//   name AS player,
	//   age AS level,
	//   job AS role
	// FROM
	//   users
}

func ExampleSprintl_Gen() {
	genYears := func(f LineFormatter) (string, bool) {
		line := f.Fmt(2020 + f.Index())
		return line, true
	}

	sql := Lines(
		"SELECT",
		"  %d", // Line 2
		"FROM",
		"  sales_data",
	).
		Gen(2, 5, genYears).
		Join(",").
		String()

	fmt.Println(sql)
	// Output:
	// SELECT
	//   2020,
	//   2021,
	//   2022,
	//   2023,
	//   2024
	// FROM
	//   sales_data
}

func ExampleSprintl_Join() {
	sql := Lines(
		"SELECT",
		"  %s", // Line 2
		"FROM",
		"  users",
	).
		Rep(2, "name", "age", "job").
		Join(",").
		String()

	fmt.Println(sql)
	// Output:
	// SELECT
	//   name,
	//   age,
	//   job
	// FROM
	//   users
}

func ExampleSprintl_Marry() {
	filters := []any{
		"name",
		"age",
		"job",
	}

	sql := Lines(
		"SELECT",
		"  name,",
		"  age,",
		"  job",
		"FROM",
		"  users",
		"WHERE",
		"%s = ?", // Line 8
	).
		Rep(8, filters...).
		Marry("  ", "AND ").
		String()

	fmt.Println(sql)
	// Output:
	// SELECT
	//   name,
	//   age,
	//   job
	// FROM
	//   users
	// WHERE
	//   name = ?
	//   AND age = ?
	//   AND job = ?
}

func ExampleSprintl_TrimSpace() {
	sql := Lines(
		"",
		"SELECT",
		"  name",
		"FROM",
		"  users",
		"",
	).
		TrimSpace().
		String()

	fmt.Println(sql)
	// Output:
	// SELECT
	//   name
	// FROM
	//   users
}

func ExampleSprintl_TrimLines() {
	sql := Lines(
		"  SELECT  ",
		"  name  ",
		"  FROM  ",
		"  users  ",
	).
		TrimLines().
		String()

	fmt.Println(sql)
	// Output:
	// SELECT
	// name
	// FROM
	// users
}

func ExampleSprintl_PruneLines() {
	sql := Lines(
		"   ",
		"SELECT",
		"  name",
		"\f",
		"		",
		"FROM",
		"  users",
		"\r\n",
	).
		PruneLines().
		String()

	fmt.Println(sql)
	// Output:
	// SELECT
	//   name
	// FROM
	//   users
}
