package nidoking

import (
	"fmt"
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

func ExampleNeedleInHaystack_Replace() {
	haystack := `
		alice, bob, charlie,
		bob, alice, charlie,
		alice, charlie, bob,
		alice, bob, charlie
	`

	var nih NeedleInHaystack
	var rep Replacement

	nih = FindNeedle(haystack, "bob", 0)
	for nih != (NeedleInHaystack{}) {
		rep = nih.Replace("dave")
		nih = rep.FindNext()
	}

	fmt.Print(trimLines(rep.Haystack))
	// Output:
	// alice, dave, charlie,
	// dave, alice, charlie,
	// alice, charlie, dave,
	// alice, dave, charlie
}

func ExampleNeedleInHaystack_ReplaceRepeat() {
	haystack := `
		alice,
		bob,
		charlie
	`

	newNames := []string{
		"jen",
		"frank",
		"ruby",
	}

	var nih NeedleInHaystack
	var rep Replacement

	nih = FindNeedle(haystack, "bob", 0)
	rep = nih.ReplaceRepeat(newNames)

	fmt.Print(trimLines(rep.Haystack))
	// Output:
	// alice,
	// jen,
	// frank,
	// ruby,
	// charlie
}
