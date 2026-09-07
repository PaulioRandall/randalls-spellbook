package sprintl

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/PaulioRandall/randalls-spellbook/pkg/nidoran"
)

type stringer interface {
	String() string
}

type mapValue struct {
	value any
	delim string
}

// Sprintl is the core type. It's public so you can pass it
// around if need, but it's intended for quick local use
// through method chaining. Objects only created via the
// [Lines] and [Split] functions.
type Sprintl struct {
	priorLineNum int
	trim         bool
	trimLines    bool
	prune        bool
	lines        []string
	formatters   map[int]LineFormatter
	// Mappings
	priorMapName string
	mappings     map[string]mapValue
}

// Lines returns a new [Sprintl] object for formatting the
// passed lines.
func Lines(lines ...string) *Sprintl {
	return &Sprintl{
		lines:      lines,
		formatters: map[int]LineFormatter{},
		mappings:   map[string]mapValue{},
	}
}

// Split returns a new [Sprintl] object for formatting the
// passed lines by splitting the passed string into lines.
func Split(s string) *Sprintl {
	return &Sprintl{
		lines:      strings.Split(s, "\n"),
		formatters: map[int]LineFormatter{},
		mappings:   map[string]mapValue{},
	}
}

// Map maps a value to a name which is referenceable within
// the template string.
//
// If the value is a struct then only its exported fields
// can be referenced. If the value is an array or slice
// then it can only be used in loop blocks. When replacing
// a value it is first checked for a String function. If
// it doesn't have one it will be stringified via
// fmt.Sprintf("%v", value).
func (s *Sprintl) Map(name string, value any) *Sprintl {
	s.mappings[name] = mapValue{
		value: value,
	}
	s.priorMapName = name
	return s
}

// MapJoin sets the join delimiter for the last mapping
// set via the [Sprintl.Map] function. If the function has
// yet to be called then panic ensues.
func (s *Sprintl) MapJoin(delim string) *Sprintl {
	if s.priorMapName == "" {
		panic("Sprintl.Map must be called before Sprintl.MapJoin")
	}

	mv := s.mappings[s.priorMapName]
	mv.delim = delim
	s.mappings[s.priorMapName] = mv
	return s
}

// Fmt format the specified line using the passed
// values.
func (s *Sprintl) Fmt(
	lineNum int,
	values ...any,
) *Sprintl {
	s.formatters[lineNum] = LineFormatter{
		typ:      "Fmt",
		values:   values,
		template: s.lines[lineNum-1],
	}
	s.priorLineNum = lineNum
	return s
}

// Dup duplicates the line with the set of args using count
// as the number of iterations. Passing 0 as the count
// removes the line completely.
func (s *Sprintl) Dup(
	lineNum int,
	count int,
	args ...any,
) *Sprintl {
	s.formatters[lineNum] = LineFormatter{
		typ:      "Dup",
		max:      count,
		values:   args,
		template: s.lines[lineNum-1],
	}
	s.priorLineNum = lineNum
	return s
}

// Rep repeats the line for each value, or array of args,
// in values. Passing an empty values array removes the
// line completely.
func (s *Sprintl) Rep(
	lineNum int,
	values ...any,
) *Sprintl {
	s.formatters[lineNum] = LineFormatter{
		typ:      "Rep",
		values:   values,
		template: s.lines[lineNum-1],
	}
	s.priorLineNum = lineNum
	return s
}

// Gen calls generator repeatedly until either it returns
// false as it second return value (notEmpty) or when the
// max iterations is reached. If notEmpty is false on the
// first call, the line is removed completely.
func (s *Sprintl) Gen(
	lineNum int,
	max int,
	generator LineGenerator,
) *Sprintl {
	s.formatters[lineNum] = LineFormatter{
		typ:       "Gen",
		max:       max,
		generator: generator,
		template:  s.lines[lineNum-1],
	}
	s.priorLineNum = lineNum
	return s
}

// Join may be called after [Sprintl.Dup], [Sprintl.Rep],
// or [Sprintl.Gen] to apply a delimiter to each line
// generated, except the last. This enables delimiter
// separated values such as comma "," separated lists. Join
// has no effect if called after [Sprintl.Fmt].
func (s *Sprintl) Join(delim string) *Sprintl {
	if s.priorLineNum == 0 {
		return s
	}

	lf, ok := s.formatters[s.priorLineNum]
	if ok {
		lf.delim = delim
		lf.prefixDelim = false
		lf.prefix = ""
		s.formatters[s.priorLineNum] = lf
	}

	return s
}

// Marry is the same as [Sprintl.Join] except the delimiter
// is applied to the start of all lines except the first.
// The first argument is a prefix which is applied to all
// lines, including the first, after the delimiter has been
// applied. This makes indented lines easy with prefixed
// delimiters.
func (s *Sprintl) Marry(prefix, delim string) *Sprintl {
	if s.priorLineNum == 0 {
		return s
	}

	lf, ok := s.formatters[s.priorLineNum]
	if ok {
		lf.delim = delim
		lf.prefixDelim = true
		lf.prefix = prefix
		s.formatters[s.priorLineNum] = lf
	}

	return s
}

// Trim may be called after [Sprintl.Fmt], [Sprintl.Dup],
// [Sprintl.Rep], or [Sprintl.Gen] to remove whitespace
// from the start and end of each line.
func (s *Sprintl) Trim() *Sprintl {
	if s.priorLineNum == 0 {
		return s
	}

	lf, ok := s.formatters[s.priorLineNum]
	if ok {
		lf.trim = true
		s.formatters[s.priorLineNum] = lf
	}

	return s
}

// TrimSpace will trim whitespace from the start and end of
// the final string. Applied after all formatters.
func (s *Sprintl) TrimSpace() *Sprintl {
	s.trim = true
	return s
}

// TrimLines will trim whitespace from every line. Applied
// after all formatters.
func (s *Sprintl) TrimLines() *Sprintl {
	s.trimLines = true
	return s
}

// PruneLines will remove all lines that are empty or only
// contain whitespace after formatting. Applied after all
// formatters.
func (s *Sprintl) PruneLines() *Sprintl {
	s.prune = true
	return s
}

// String compiles each line, joins them together, and
// returns the result.
func (s *Sprintl) String() string {
	lines := applyLineFormatters(s.lines, s.formatters)

	if s.prune {
		lines = prune(lines)
	}

	if s.trimLines {
		trimLines(lines)
	}

	str := strings.Join(lines, "\n")
	str = s.applyMappings(str)

	if s.trim {
		return strings.TrimSpace(str)
	}

	return str
}

func applyLineFormatters(
	lines []string,
	formatters map[int]LineFormatter,
) []string {
	result := []string{}

	for i, line := range lines {
		lineNum := i + 1

		form, ok := formatters[lineNum]
		if !ok {
			result = append(result, line)
			continue
		}

		result = append(result, form.apply()...)
	}

	return result
}

func (s *Sprintl) applyMappings(text string) string {
	for name, mapValue := range s.mappings {

		text = applyMapping(text, name, mapValue)
	}
	return text
}

func applyMapping(
	text string,
	name string,
	mv mapValue,
) string {
	if IsArrayOrSlice(mv.value) {
		return applyLineMapping(text, name, mv)
	}
	return applyValueMapping(text, name, mv)
}

func createSearchString(name string) string {
	return "{{" + name + "}}"
}

func applyValueMapping(
	text, name string,
	mv mapValue,
) string {
	needle := createSearchString(name)
	//value := stringifyValue(mv.value)
	nih := nidoran.Find(text, needle, 0)

	for nih != (nidoran.NeedleInHaystack{}) {
		//nih = nih.ReplaceFindNext(value)
	}

	return text
}

func applyLineMapping(
	text, name string,
	mv mapValue,
) string {
	needle := createSearchString(name)
	//values := stringifySliceValues(mv)
	nih := nidoran.Find(text, needle, 0)

	// TODO

	for nih != (nidoran.NeedleInHaystack{}) {
		//nih = nih.ReplaceLineFindNext(line)
	}

	return ""
}

func IsArrayOrSlice(v any) bool {
	k := reflect.TypeOf(v).Kind()
	return k == reflect.Array || k == reflect.Slice
}

func stringifySliceValues(mv mapValue) []string {
	refValue := reflect.ValueOf(mv.value)
	result := make([]string, refValue.Len())

	for i := 0; i < refValue.Len(); i++ {
		v := refValue.Index(i).Interface()
		result[i] = stringifyValue(v)
	}

	return result
}

func toAnySlice(values any) []any {
	refValue := reflect.ValueOf(values)
	result := make([]any, refValue.Len())

	for i := 0; i < refValue.Len(); i++ {
		result[i] = refValue.Index(i).Interface()
	}

	return result
}

func stringifyValue(value any) string {
	if obj, ok := value.(stringer); ok {
		return obj.String()
	}

	return fmt.Sprintf("%v", value)
}

func trimLines(lines []string) {
	for i, s := range lines {
		lines[i] = strings.TrimSpace(s)
	}
}

func prune(lines []string) []string {
	result := []string{}

	for _, s := range lines {
		if strings.TrimSpace(s) != "" {
			result = append(result, s)
		}
	}

	return result
}
