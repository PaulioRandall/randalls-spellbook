package nidorino

import (
	"fmt"
	"strings"

	"github.com/PaulioRandall/randalls-spellbook/pkg/nidoran"
)

var emptyNih = nidoran.NeedleInHaystack{}

type stringer interface {
	String() string
}

// Formatter is returned by [Given] and [Lines] and
// provides functions for populating mustache placeholder
// tokens, e.g. '{{token}}', in text. All formatting and
// print functions return the formatter for method
// chaining.
type Formatter struct {
	text string
}

// Given returns a new Formatter for formatting the given
// template string.
func Given(template string) *Formatter {
	return &Formatter{
		text: template,
	}
}

// Given returns a new Formatter for formatting the given
// template string, provided as a series of lines.
func Lines(lines ...string) *Formatter {
	return &Formatter{
		text: strings.Join(lines, "\n"),
	}
}

// Print prints the text to terminal primarily for the
// purposes of debugging. Formatting is performed as
// formatting functions are called so this function can be
// called after each formatting function if need be. The
// Formatter is returned for method chaining so it can be
// used mid-chain.
func (form *Formatter) Print() *Formatter {
	fmt.Println(form.text)
	return form
}

// PrintNamed does the same as [Formatter.Println] except
// delimiter lines are printed before and after the text
// containing the passed label so it's clearer which output
// corrisponds to which printing.
func (form *Formatter) PrintNamed(label string) *Formatter {
	fmt.Println(">>>" + label + ">>>")
	fmt.Println(form.text)
	fmt.Println("<<<" + label + "<<<")
	return form
}

// Fmt populates every token named key with value.
func (form *Formatter) Fmt(
	key string,
	value any,
) *Formatter {
	nih := nidoran.Find(form.text, token(key), 0)
	str := stringifyValue(value)

	for nih != (emptyNih) {
		rep := nih.ReplaceInline(str)
		form.text = rep.Haystack
		nih = rep.FindNext()
	}

	return form
}

// Join replaces every line contining token named key with
// a set of lines where the token is replaced values. Each
// line is suffixed with delim, except the last.
func (form *Formatter) Join[T any](
	key, delim string,
	values ...T,
) *Formatter {
	nih := nidoran.Find(form.text, token(key), 0)
	strs := stringifyValues(values)

	for nih != (emptyNih) {
		rep := nih.ReplaceJoin(strs, ",")
		form.text = rep.Haystack
		nih = rep.FindNext()
	}

	return form
}

// String returns the text. All formatting functions called
// prior will have been applied to the string.
func (form *Formatter) String() string {
	return form.text
}

func token(key string) string {
	return "{{" + key + "}}"
}

func stringifyValue[T any](value T) string {
	if st, ok := any(value).(stringer); ok {
		return st.String()
	}

	return fmt.Sprintf("%v", value)
}

func stringifyValues[T any](values []T) []string {
	strs := make([]string, len(values), len(values))

	for i, v := range values {
		strs[i] = stringifyValue(v)
	}

	return strs
}
