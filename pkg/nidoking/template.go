package nidoking

import (
	"fmt"
	"slices"
	"strings"
)

func token(key string) string {
	return "{{" + key + "}}"
}

// Template is returned by [Given] and [Lines] and
// provides functions for populating mustache placeholder
// tokens, e.g. '{{key}}', in text. All formatting and
// print functions return a reference to the Template for
// method chaining.
type Template struct {
	text string
}

// Given returns a new Template for formatting the given
// template string.
func Given(template string) *Template {
	return &Template{
		text: template,
	}
}

// Given returns a new Template with the template string
// produced from the given lines.
func Lines(lines ...string) *Template {
	return &Template{
		text: strings.Join(lines, "\n"),
	}
}

func (tmpl *Template) replace[T any, S string | []string](
	key string,
	value T,
	stringify func(T) S,
	f func(nih NeedleInHaystack, s S) Replacement,
) *Template {
	var nih NeedleInHaystack
	var rep Replacement

	v := stringify(value)
	nih = Find(tmpl.text, token(key), 0)

	for nih != (NeedleInHaystack{}) {
		rep = f(nih, v)
		nih = rep.FindNext()
	}

	tmpl.text = rep.Haystack
	return tmpl
}

// Print prints the text to terminal primarily for the
// purposes of debugging. Formatting is performed as
// formatting functions are called so this function can be
// called after each format to see how the text evolves.
func (tmpl *Template) Print() *Template {
	fmt.Println(tmpl.text)
	return tmpl
}

// PrintNamed does the same as [Template.Print] except
// delimiter lines are printed before and after the text,
// containing the passed label, so it's clearer which
// output corrisponds to which print statement.
func (tmpl *Template) PrintNamed(label string) *Template {
	fmt.Println(">>>" + label + ">>>")
	fmt.Println(tmpl.text)
	fmt.Println("<<<" + label + "<<<")
	return tmpl
}

// Fmt populates every token named key with value.
func (tmpl *Template) Fmt(
	key string,
	value any,
) *Template {
	f := func(nih NeedleInHaystack, s string) Replacement {
		return nih.ReplaceInline(s)
	}

	return tmpl.replace(
		key,
		value,
		stringifyValue,
		f,
	)
}

// FmtJoin replaces every token named key with the list of
// values. Each instance of value is suffixed with delim,
// except the last.
func (tmpl *Template) FmtJoin[T any](
	key, delim string,
	values ...T,
) *Template {
	f := func(nih NeedleInHaystack, s []string) Replacement {
		return nih.ReplaceInlineJoin(s, delim)
	}

	return tmpl.replace(
		key,
		values,
		stringifyValues,
		f,
	)
}

// FmtRepeat replaces every token named key with value
// repeated n times. Each instance of value is suffixed
// with delim, except the last.
func (tmpl *Template) FmtRepeat[T any](
	key, delim string,
	value T,
	n int,
) *Template {
	f := func(nih NeedleInHaystack, s string) Replacement {
		return nih.ReplaceInlineRepeat(s, n, delim)
	}

	return tmpl.replace(
		key,
		value,
		stringifyValue,
		f,
	)
}

// Join replaces every line containing a token named key
// with a set of lines where the token is replaced by
// values. Each line is suffixed with delim, except the
// last.
func (tmpl *Template) Join[T any](
	key, delim string,
	values ...T,
) *Template {
	f := func(nih NeedleInHaystack, s []string) Replacement {
		return nih.ReplaceJoin(s, delim)
	}

	return tmpl.replace(
		key,
		values,
		stringifyValues,
		f,
	)
}

// Repeat replaces every line containing a token named key
// with n lines where the token is replaced by value in
// each. Each line is suffixed with delim, except the last.
func (tmpl *Template) Repeat[T any](
	key, delim string,
	value T,
	n int,
) *Template {
	f := func(nih NeedleInHaystack, s string) Replacement {
		return nih.ReplaceRepeat(s, n, delim)
	}

	return tmpl.replace(
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
func (tmpl *Template) Map[T any](
	key string,
	max int,
	gen func(i int) (T, bool),
) *Template {
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

	return tmpl.replace(
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
func (tmpl *Template) CopyLines(
	start, end, at, amount int,
) *Template {
	lines := strings.Split(tmpl.text, "\n")
	copies := lines[start:end]

	for i := 0; i < amount; i++ {
		lines = slices.Insert(
			lines,
			at,
			copies...,
		)
	}

	tmpl.text = strings.Join(lines, "\n")
	return tmpl
}

// String returns the text with all formatting applied.
func (tmpl *Template) String() string {
	return tmpl.text
}
