package wizzard

import (
	"fmt"
)

func Example() {
	abc := func(i int, s string) (int, error) {
		return len(s) + i, nil
	}

	spellAbc := NewSpell("abc", abc)

	result, e := spellAbc.Invoke(100, "hello")

	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(result)
	}
	// Output:
	// 105
}
