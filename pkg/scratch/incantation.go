package scratch

import (
	"fmt"
)

// Incantation is a single step in casting a spell.
// Incantations are weaved together by converting them into
// Thunks.
type Incantation[C, T, R any] func(C, T) (R, error)

// Spell is a Thunk encapsulating an Incantation call,
// returning the output.
type Spell[R any] func() (R, error)

// WeaveInit returns a Spell for the given input and
// incantation. It is called to create the first part of
// a spell.
func WeaveInit[C, T, R any](
	ctx C,
	input T,
	incant Incantation[C, T, R],
) Spell[R] {
	return func() (R, error) {
		return incant(ctx, input)
	}
}

// Weave appends to a spell by returning a new Spell that
// encapsulates an existing Spell (the spell so far) and
// the next Incantation.
func Weave[C, T, R any](
	ctx C,
	spell Spell[T],
	incant Incantation[C, T, R],
) Spell[R] {
	return func() (R, error) {
		result, e := spell()

		if e != nil {
			var empty R
			return empty, e
		}

		return incant(ctx, result)
	}
}

// EXAMPLE

type numberNames = map[int]string

func incantationExample() {
	ctx := numberNames{
		1: "one",
		2: "two",
		3: "three",
		4: "four",
		5: "five",
		6: "six",
		7: "seven",
		8: "eight",
		9: "nine",
	}

	input := 2
	spellPart1 := WeaveInit(ctx, input, findNumberName)
	spellPart2 := Weave(ctx, spellPart1, countRunesInString)
	finalSpell := Weave(ctx, spellPart2, findNameForNumber)

	var namedLength string
	var e error

	namedLength, e = finalSpell()

	if e != nil {
		println(e)
	} else {
		msg := fmt.Sprintf("The number %d as a word has %s letters", input, namedLength)
		println(msg)
	}
}

func findNumberName(ctx numberNames, i int) (string, error) {
	result := ctx[i]
	if result == "" {
		return "", fmt.Errorf("Name not found for number %d", i)
	}
	return result, nil
}

func countRunesInString(_ numberNames, s string) (int, error) {
	return len([]rune(s)), nil
}

func findNameForNumber(ctx numberNames, i int) (string, error) {
	for number, name := range ctx {
		if number == i {
			return name, nil
		}
	}
	return "", fmt.Errorf("Missing entry for the number %d", i)
}
