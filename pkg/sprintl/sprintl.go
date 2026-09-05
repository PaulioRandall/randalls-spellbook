package sprintl

import (
	"strings"
)

// LineFormatter is passed to the [Sprintl.Gen] function to
// enable customised formatting.
type LineFormatter func(
	formatter LineFormat,
) (
	result string,
	notEmpty bool,
)

// Sprintl is the core type. It's public so you can pass it
// around if need, but it's intended for quick local use
// through method chaining. Call [Lines] to create an
// object.
type Sprintl struct {
	priorLineNum int
	lines        []string
	formatters   map[int]LineFormat
}

// Lines returns a new [Sprintl] object for formatting the
// passed lines.
func Lines(lines ...string) *Sprintl {
	return &Sprintl{
		lines:      lines,
		formatters: map[int]LineFormat{},
	}
}

// Fmt formatters the specified line using the passed
// values.
func (s *Sprintl) Fmt(
	lineNum int,
	values ...any,
) *Sprintl {
	s.formatters[lineNum] = LineFormat{
		typ:      "Fmt",
		values:   values,
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
	s.formatters[lineNum] = LineFormat{
		typ:      "Rep",
		values:   values,
		template: s.lines[lineNum-1],
	}
	s.priorLineNum = lineNum
	return s
}

// Gen calls generator repeatedly until either
// [LineFormatter] notEmpty return value is false or the
// max iterations is reached. If notEmpty is false on the
// first call, the line is removed completely.
func (s *Sprintl) Gen(
	lineNum int,
	max int,
	generator LineFormatter,
) *Sprintl {
	s.formatters[lineNum] = LineFormat{
		typ:       "Gen",
		max:       max,
		generator: generator,
		template:  s.lines[lineNum-1],
	}
	s.priorLineNum = lineNum
	return s
}

// Join may be called after [Sprintl.Rep] or [Sprintl.Gen]
// to apply a delimiter to each line generated, except the
// last line. This enables delimiter separated lists. Most
// commonly a comma "," is used as the delim value. Join
// has no effect if called after [Sprintl.Fmt] or other
// functions.
func (s *Sprintl) Join(delim string) *Sprintl {
	if s.priorLineNum == 0 {
		return s
	}

	f, ok := s.formatters[s.priorLineNum]
	if ok {
		f.delim = delim
		f.prefixDelim = false
		s.formatters[s.priorLineNum] = f
	}

	return s
}

// Marry is the same as [Sprintl.Join] except the delimiter
// is applied to the start of all lines except the first.
func (s *Sprintl) Marry(delim string) *Sprintl {
	if s.priorLineNum == 0 {
		return s
	}

	f, ok := s.formatters[s.priorLineNum]
	if ok {
		f.delim = delim
		f.prefixDelim = true
		s.formatters[s.priorLineNum] = f
	}

	return s
}

// Prefix adds a prefix to each line. This is primarily
// designed to be used along with [Sprintl.Marry] to insert
// line indents. It is applied to every line in the
// repetition after [Sprintl.Marry] is applied.
func (s *Sprintl) Prefix(prefix string) *Sprintl {
	if s.priorLineNum == 0 {
		return s
	}

	f, ok := s.formatters[s.priorLineNum]
	if ok {
		f.prefix = prefix
		s.formatters[s.priorLineNum] = f
	}

	return s
}

// String compiles each line, joins them together, and
// returns the result.
func (s *Sprintl) String() string {
	lines := []string{}

	for i, line := range s.lines {
		lineNum := i + 1

		form, ok := s.formatters[lineNum]
		if !ok {
			lines = append(lines, line)
			continue
		}

		lines = append(lines, form.apply()...)
	}

	return strings.Join(lines, "\n")
}
