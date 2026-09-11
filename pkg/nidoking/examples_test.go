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

func ExampleFind() {
	haystack := `
	alice,
	😁bob,
	charlie
`

	nih := Find(haystack, "bob", 0)
	fmt.Println(nih.String())
	// Output:
	// NeedleInHaystack{
	//   LineStart: 9,
	//   Start: 14,
	//   End: 17,
	//   LineEnd: 18,
	//   Mode: "string",
	//   Needle: "bob",
	//   Pattern: "bob",
	//   Haystack: "\n	alice,\n	😁bob,\n	charlie\n",
	//   IsMatch(): true,
	//   LineIndex(): 2,
	//   LineNumber(): 3,
	//   InlineStart(): 5,
	//   InlineEnd(): 8,
	//   RuneLineStart(): 9,
	//   RuneStart(): 11,
	//   RuneEnd(): 14,
	//   RuneLineEnd(): 15,
	//   RuneInlineStart(): 2,
	//   RuneInlineEnd(): 5,
	// }
}

func ExampleMatch() {
	haystack := `
	alice,
	😁bob,
	charlie
`

	nih := Match(haystack, "ch.*lie", 0)
	fmt.Println(nih.String())
	// Output:
	// NeedleInHaystack{
	//   LineStart: 19,
	//   Start: 20,
	//   End: 27,
	//   LineEnd: 27,
	//   Mode: "regexp",
	//   Needle: "charlie",
	//   Pattern: "ch.*lie",
	//   Haystack: "\n	alice,\n	😁bob,\n	charlie\n",
	//   IsMatch(): true,
	//   LineIndex(): 3,
	//   LineNumber(): 4,
	//   InlineStart(): 1,
	//   InlineEnd(): 8,
	//   RuneLineStart(): 16,
	//   RuneStart(): 17,
	//   RuneEnd(): 24,
	//   RuneLineEnd(): 24,
	//   RuneInlineStart(): 1,
	//   RuneInlineEnd(): 8,
	// }
}

func ExampleNeedleInHaystack_FindNext() {
	// Note: the first and last lines are whitespace only!
	haystack := `
	alice, bob, charlie,
	bob, alice, charlie,
	alice, charlie, bob
	alice, bob, charlie
`

	var matches []NeedleInHaystack

	nih := Find(haystack, "bob", 0)
	for nih != (NeedleInHaystack{}) {
		matches = append(matches, nih)
		nih = nih.FindNext()
	}

	for i, nih := range matches {
		fmt.Printf(
			"[%d] line: %d [%d:%d]\n",
			i,
			nih.LineNumber(),
			nih.InlineStart(),
			nih.InlineEnd(),
		)
	}
	// Output:
	// [0] line: 2 [8:11]
	// [1] line: 3 [1:4]
	// [2] line: 4 [17:20]
	// [3] line: 5 [8:11]
}

func ExampleNeedleInHaystack_ReplaceInline() {
	haystack := `
		alice, bob, charlie,
		bob, alice, charlie,
		alice, charlie, bob,
		alice, bob, charlie
	`

	var nih NeedleInHaystack
	var rep Replacement

	nih = Find(haystack, "bob", 0)
	for nih != (NeedleInHaystack{}) {
		rep = nih.ReplaceInline("dave")
		nih = rep.FindNext()
	}

	fmt.Print(trimLines(rep.Haystack))
	// Output:
	// alice, dave, charlie,
	// dave, alice, charlie,
	// alice, charlie, dave,
	// alice, dave, charlie
}

func ExampleNeedleInHaystack_ReplaceJoin() {
	haystack := joinLines(
		"SELECT",
		"	{{columns}}",
		"FROM",
		"	players",
	)

	columns := []string{
		"name",
		"level",
		"role",
	}

	var nih NeedleInHaystack
	var rep Replacement

	nih = Find(haystack, "{{columns}}", 0)
	rep = nih.ReplaceJoin(columns, ",")

	fmt.Print(rep.Haystack)
	// Output:
	// SELECT
	// 	name,
	// 	level,
	// 	role
	// FROM
	// 	players
}
