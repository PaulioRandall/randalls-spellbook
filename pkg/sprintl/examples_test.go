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
		Rep(2, "name", "age", "job").
		Join(",").
		Fmt(4, "users").
		Gen(6, 2, func(f LineFormat) (string, bool) {
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
		Rep(2, "name", "age", "height").
		Join(",").
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
	props := []any{
		"name",
		"age",
		"job",
	}

	md := Lines(
		"# Player Properties",
		"",
		"%s", // Line 3
	).
		Rep(3, props...).
		Marry("and ").
		String()

	fmt.Println(md)
	// Output:
	// # Player Properties
	//
	// name
	// and age
	// and job
}

func ExampleSprintl_Prefix() {
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
		Marry("AND ").
		Prefix("  ").
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
