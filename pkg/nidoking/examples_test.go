package nidoking

import (
	"fmt"
	"strings"
)

func Example() {
	haystack := `
	alice,
	😁bob,
	charlie
`

	nih := FindNeedle(haystack, "bob", 0)
	fmt.Println(nih.String())
	// Output:
	// NeedleInHaystack{
	//   LineIndex: 2,
	//   LineStart: 9,
	//   LineEnd: 18,
	//   Start: 14,
	//   End: 17,
	//   Needle: 'bob',
	//   Haystack: <not printed>,
	//   LineNum(): 3,
	//   InlineStart(): 5,
	//   InlineEnd(): 8,
	//   RuneLineStart(): 9,
	//   RuneLineEnd(): 15,
	//   RuneStart(): 11,
	//   RuneEnd(): 14,
	//   RuneInlineStart(): 2,
	//   RuneInlineEnd(): 5,
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
	// [0] line: 2 [8:11]
	// [1] line: 3 [1:4]
	// [2] line: 4 [17:20]
	// [3] line: 5 [8:11]
}

func ExampleNeedleInHaystack_ReplaceFindNext() {
	haystack := `
		alice, bob, charlie,
		bob, alice, charlie,
		alice, charlie, bob,
		alice, bob, charlie
	`

	nih := FindNeedle(haystack, "bob", 0)
	for nih != (NeedleInHaystack{}) {
		haystack = nih.Replace("dave")
		nih = nih.ReplaceFindNext("dave")
	}

	fmt.Print(trimLines(haystack))
	// Output:
	// alice, dave, charlie,
	// dave, alice, charlie,
	// alice, charlie, dave,
	// alice, dave, charlie
}

// Shows how the haystack differs between the result
// objects of ReplaceFindNext and FindNext.
func ExampleNeedleInHaystack_ReplaceFindNext_diff() {
	getTrimmedFirstLine := func(haystack string) string {
		haystack = strings.TrimSpace(haystack)
		return strings.Split(haystack, "\n")[0]
	}

	haystack := `
	alice, bob, charlie
	bob, alice, charlie
`

	nih := FindNeedle(haystack, "bob", 0)

	a := nih.ReplaceFindNext("dave").Haystack
	b := nih.FindNext().Haystack

	a = getTrimmedFirstLine(a)
	b = getTrimmedFirstLine(b)

	fmt.Printf("'%s' (ReplaceFindNext)\n", a)
	fmt.Printf("'%s'  (FindNext)\n", b)
	// Output:
	// 'alice, dave, charlie' (ReplaceFindNext)
	// 'alice, bob, charlie'  (FindNext)
}
