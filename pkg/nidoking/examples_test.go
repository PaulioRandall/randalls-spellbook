package nidoking

import (
	"fmt"
	"strings"
)

func Example() {
	haystack := `
		alice,
		bob,
		charlie
	`

	nih := FindNeedle(haystack, "bob", 0)
	fmt.Println(nih.String())

	// Output:
	// NeedleInHaystack{
	//   LineIndex: 2,
	//   LineStart: 10,
	//   LineEnd: 16,
	//   Start: 12,
	//   End: 15,
	//   Needle: bob,
	//   Haystack: <not printed>,
	//   LineNum(): 3,
	//   InlineStart(): 2,
	//   InlineEnd(): 5,
	// }
}

func ExampleNeedleInHaystack_FindNext() {
	// The first and last lines are whitespace only.
	// I missed it too when writing this example!
	haystack := `
		alice, bob, charlie,
		bob, alice, charlie,
		alice, charlie, bob
		alice, bob, charlie
	`

	var matches []NeedleInHaystack

	nih := FindNeedle(haystack, "bob", 0)
	for nih != (NeedleInHaystack{}) {
		matches = append(matches, nih)
		nih = nih.FindNext()
	}

	fmt.Printf(
		"Found %d instances of 'bob'\n",
		len(matches),
	)

	for i, nih := range matches {
		fmt.Printf(
			"[%d] line: %d [%d:%d]\n",
			i,
			nih.LineNum(),
			nih.InlineStart(),
			nih.InlineEnd(),
		)
	}

	// Output:
	// Found 4 instances of 'bob'
	// [0] line: 2 [9:12]
	// [1] line: 3 [2:5]
	// [2] line: 4 [18:21]
	// [3] line: 5 [9:12]
}

func ExampleNeedleInHaystack_ReplaceNeedleFindNext() {
	// The first and last lines are whitespace only.
	// I missed it too when writing this example!
	haystack := `
		alice, bob, charlie,
		bob, alice, charlie,
		alice, charlie, bob,
		alice, bob, charlie
	`

	nih := FindNeedle(haystack, "bob", 0)

	for nih != (NeedleInHaystack{}) {
		haystack = nih.ReplaceNeedle("dave")

		// Calling nih.FindNext will not produce the same
		// result because it won't use the updated haystack!
		nih = nih.ReplaceNeedleFindNext("dave")
	}

	fmt.Print(trimLines(haystack))

	// Output:
	// alice, dave, charlie,
	// dave, alice, charlie,
	// alice, charlie, dave,
	// alice, dave, charlie
}

func ExampleNeedleInHaystack_ReplaceNeedleFindNext_diff() {
	// The first and last lines are whitespace only.
	// I missed it too when writing this example!
	haystack := `
		alice, bob, charlie
		bob, alice, charlie
	`

	nih := FindNeedle(haystack, "bob", 0)

	getTrimmedFirstLine := func(haystack string) string {
		haystack = strings.TrimSpace(haystack)
		return strings.Split(haystack, "\n")[0]
	}

	a := nih.ReplaceNeedle("dave") // Returns new haystack
	b := nih.ReplaceNeedleFindNext("dave").Haystack
	c := nih.FindNext().Haystack

	a = getTrimmedFirstLine(a)
	b = getTrimmedFirstLine(b)
	c = getTrimmedFirstLine(c)

	fmt.Printf("'%s' (ReplaceNeedle)\n", a)
	fmt.Printf("'%s' (ReplaceNeedleFindNext)\n", b)
	fmt.Printf("'%s'  (FindNext)\n", c)

	// Output:
	// 'alice, dave, charlie' (ReplaceNeedle)
	// 'alice, dave, charlie' (ReplaceNeedleFindNext)
	// 'alice, bob, charlie'  (FindNext)
}
