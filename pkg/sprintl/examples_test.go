package sprintl

import (
	"fmt"
)

func Example() {
	conditions := []string{
		"name = ?",
		"age = ?",
		"job = ?",
	}

	sql := Lines(
		"SELECT",
		"  %s", // Line 2
		"FROM",
		"  %s", // Line 4
		"WHERE",
		"  %s", // Line 6
	).
		Rep(2, ",", "name", "age", "job").
		Fmt(4, "users").
		Gen(6, "", 2, func(f LineFormat) (string, bool) {
			// The 2 above is the max number of repetitions.
			cond := conditions[f.Index()]

			if f.Index() > 0 {
				cond = "AND " + cond
			}

			// Returning false will end repetition.
			return f.Fmt(cond), true
		})

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

func ExampleSprintl_Rep() {
	sql := Lines(
		"SELECT",
		"  %s", // Line 2
		"FROM",
		"  users",
		"WHERE",
		"  %s = ?", // Line 6
	).
		Rep(2, ",", "name", "age", "height").
		Fmt(6, "name").
		String()

	fmt.Println(sql)
	// Output:
	// SELECT
	//   name,
	//   age,
	//   height
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
		Rep(2, ",",
			[]any{"name", "player"},
			[]any{"age", "level"},
			[]any{"job", "role"},
		).
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
	genYears := func(f LineFormat) (string, bool) {
		line := f.Fmt(2020 + f.Index())
		return line, true
	}

	sql := Lines(
		"SELECT",
		"  %d", // Line 2
		"FROM",
		"  sales_data",
	).
		Gen(2, ",", 5, genYears).
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
