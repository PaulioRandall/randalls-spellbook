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

// String returns the text with all formatting applied.
func (tmpl *Template) String() string {
	return tmpl.text
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

// FmtObject searches for tokens with the pattern
// {{key.field}} and replaces it with the value named
// 'field' in object, or the value of calling the method
// named 'field'. The method must accept no arguments and
// return a single value.
func (tmpl *Template) FmtObject(
	key string,
	object any,
) *Template {
	panic("TODO: Template.FmtObject")
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

// Objects searches for lines with tokens containing the
// pattern {{key.field}} and repeats each line with the
// value named 'field' in each object, or the value of
// calling the method named 'field'. The method must
// accept no arguments and return a single value. Each line
// is suffixed with delim, except the last.
func (tmpl *Template) Objects[T any](
	key, delim string,
	objects ...T,
) *Template {
	panic("TODO: Template.Objects")
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

// RemoveLines removes all lines between startKey and
// endKey. The lines containing the tokens are also
// removed.
func (tmpl *Template) RemoveLines(
	startKey, endKey string,
) *Template {
	return tmpl.RepeatLines(startKey, endKey, 0)
}

// KeepLines removes the lines containining startKey and
// endKey leaving the content.
func (tmpl *Template) KeepLines(
	startKey, endKey string,
) *Template {
	return tmpl.RepeatLines(startKey, endKey, 1)
}

// RepeatLines repeats the set of lines between startKey
// and endKey (exclusive) by amount. The lines containing
// the tokens are also removed. Passing 0 as the amount
// will remove all content; same as calling
// [Template.RemoveLines].
func (tmpl *Template) RepeatLines(
	startKey, endKey string,
	amount int,
) *Template {
	start, end := tmpl.locateDelimiterLineIndexes(
		startKey,
		endKey,
	)

	lines := strings.Split(tmpl.text, "\n")
	copy := lines[start+1 : end]

	lines = append(lines[:start], lines[end+1:]...)
	if amount == 0 {
		tmpl.text = strings.Join(lines, "\n")
		return tmpl
	}

	copy = slices.Repeat(copy, amount)
	lines = slices.Insert(lines, start, copy...)
	tmpl.text = strings.Join(lines, "\n")
	return tmpl
}

func (tmpl *Template) locateDelimiterLineIndexes(
	startKey, endKey string,
) (int, int) {
	startNih := Find(tmpl.text, token(startKey), 0)
	if startNih == (NeedleInHaystack{}) {
		msg := fmt.Errorf(
			"Start token key '%s' not found",
			startKey,
		)
		panic(msg)
	}

	endNih := Find(tmpl.text, token(endKey), startNih.End)
	if endNih == (NeedleInHaystack{}) {
		msg := fmt.Errorf(
			"End token key '%s' not found",
			endKey,
		)
		panic(msg)
	}

	startLineIndex := startNih.LineIndex()
	endLineIndex := endNih.LineIndex()

	if startLineIndex == endLineIndex {
		msg := fmt.Errorf(
			"Start and end token keys, '%s' and '%s', cannot be on the same line",
			startKey,
			endKey,
		)
		panic(msg)
	}

	return startLineIndex, endLineIndex
}
