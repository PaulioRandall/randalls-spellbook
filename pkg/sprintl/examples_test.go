package sprintl

import (
	"fmt"
)

func ExampleLines() {
	var _ *Sprintl = Lines(
		"SELECT",
		"  %s", // Line 2
		"FROM",
		"  %s", // Line 4
	)
}

func ExampleSprintl_F() {
	sql := Lines(
		"SELECT",
		"  'month_' || '%d' AS %s,", // Line 2
		"  revenue",
		"FROM",
		"  sales_data",
		"WHERE",
		"  year = 2025",
	).
		F(2, 9, "sept").
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

func ExampleSprintl_R() {
	sql := Lines(
		"SELECT",
		"  %s", // Line 2
		"FROM",
		"  users",
		"WHERE",
		"  %s = ?", // Line 6
	).
		R(2, ",", "name", "age", "height").
		F(6, "name").
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

func ExampleSprintl_G() {
	genYears := func(f Formatter) (string, bool) {
		line := f.Fmt(2020 + f.Index())
		return line, true
	}

	sql := Lines(
		"SELECT",
		"  %d", // Line 2
		"FROM",
		"  sales_data",
	).
		G(2, ",", 5, genYears).
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
