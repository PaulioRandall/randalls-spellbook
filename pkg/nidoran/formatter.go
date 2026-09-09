package nidoran

import (
	"fmt"
	"slices"
	"strings"
)

func token(key string) string {
	return "{{" + key + "}}"
}

// Formatter is returned by [Given] and [Lines] and
// provides functions for populating mustache placeholder
// tokens, e.g. '{{token}}', in text. All formatting and
// print functions return the Formatter for method
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
// lines.
func Lines(lines ...string) *Formatter {
	return &Formatter{
		text: strings.Join(lines, "\n"),
	}
}

func (form *Formatter) replace[T any, S string | []string](
	key string,
	value T,
	stringify func(T) S,
	f func(nih NeedleInHaystack, s S) Replacement,
) *Formatter {
	var nih NeedleInHaystack
	var rep Replacement

	v := stringify(value)
	nih = Find(form.text, token(key), 0)

	for nih != (NeedleInHaystack{}) {
		rep = f(nih, v)
		nih = rep.FindNext()
	}

	form.text = rep.Haystack
	return form
}

// Print prints the text to terminal primarily for the
// purposes of debugging. Formatting is performed as
// formatting functions are called so this function can be
// called after each format to see how the text evolves.
func (form *Formatter) Print() *Formatter {
	fmt.Println(form.text)
	return form
}

// PrintNamed does the same as [Formatter.Print] except
// delimiter lines are printed before and after the text,
// containing the passed label, so it's clearer which
// output corrisponds to which print statement.
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
	f := func(nih NeedleInHaystack, s string) Replacement {
		return nih.ReplaceInline(s)
	}

	return form.replace(
		key,
		value,
		stringifyValue,
		f,
	)
}

// FmtJoin replaces every token named key with the list of
// values. Each instance of value is suffixed with delim,
// except the last.
func (form *Formatter) FmtJoin[T any](
	key, delim string,
	values ...T,
) *Formatter {
	f := func(nih NeedleInHaystack, s []string) Replacement {
		return nih.ReplaceInlineJoin(s, delim)
	}

	return form.replace(
		key,
		values,
		stringifyValues,
		f,
	)
}

// FmtRepeat replaces every token named key with value
// repeated n times. Each instance of value is suffixed
// with delim, except the last.
func (form *Formatter) FmtRepeat[T any](
	key, delim string,
	value T,
	n int,
) *Formatter {
	f := func(nih NeedleInHaystack, s string) Replacement {
		return nih.ReplaceInlineRepeat(s, n, delim)
	}

	return form.replace(
		key,
		value,
		stringifyValue,
		f,
	)
}

// Join replaces every line containing a token named key
// with a set of lines where the token is replaced values.
// Each line is suffixed with delim, except the last.
func (form *Formatter) Join[T any](
	key, delim string,
	values ...T,
) *Formatter {
	f := func(nih NeedleInHaystack, s []string) Replacement {
		return nih.ReplaceJoin(s, delim)
	}

	return form.replace(
		key,
		values,
		stringifyValues,
		f,
	)
}

// Repeat replaces every line containing a token named key
// with n lines where the token is replaced by value in
// each. Each line is suffixed with delim, except the last.
func (form *Formatter) Repeat[T any](
	key, delim string,
	value T,
	n int,
) *Formatter {
	f := func(nih NeedleInHaystack, s string) Replacement {
		return nih.ReplaceRepeat(s, n, delim)
	}

	return form.replace(
		key,
		value,
		stringifyValue,
		f,
	)
}

// Map replaces every line containing a token named key
// with the same line but the token replaced with the
// result of calling gen. If gen returns false then the
// replacement finishes with gen's last string value
// ignored. The maximum number of iterations is determined
// by max.
//
// If the caller knows how many iterations will be needed
// then it's recommended to set it as max and always
// return true from gen. If max is not known, set a higher
// than expected value for max, e.g. math.MaxUint8, and
// control iteration exit via gen's return bool.
func (form *Formatter) Map[T any](
	key string,
	max int,
	gen func(i int) (T, bool),
) *Formatter {
	var values []any

	for i := 0; i < max; i++ {
		v, ok := gen(i)
		if !ok {
			break
		}

		values = append(values, v)
	}

	f := func(nih NeedleInHaystack, s []string) Replacement {
		return nih.ReplaceJoin(s, "")
	}

	return form.replace(
		key,
		values,
		stringifyValues,
		f,
	)
}

// CopyLines copies the series of lines from start to end
// (exclusive) and inserts them in front of line 'at'. No
// formatting is performed and no lines are overwrtten or
// removed.
func (form *Formatter) CopyLines(
	start, end, at, amount int,
) *Formatter {
	lines := strings.Split(form.text, "\n")
	copies := lines[start:end]

	for i := 0; i < amount; i++ {
		lines = slices.Insert(
			lines,
			at,
			copies...,
		)
	}

	form.text = strings.Join(lines, "\n")
	return form
}

// String returns the text. All formatting functions called
// prior will have been applied to the string.
func (form *Formatter) String() string {
	return form.text
}
