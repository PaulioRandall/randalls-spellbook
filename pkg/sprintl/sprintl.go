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
	lines      []string
	formatters map[int]LineFormat
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
	return s
}

// Rep repeats the line for each value in values then
// appends the delim to each, except the last line. The
// line must contain a single placeholder, e.g. '%v'.
// Passing an empty values array removes the line
// completely.
func (s *Sprintl) Rep(
	lineNum int,
	delim string,
	values ...any,
) *Sprintl {
	s.formatters[lineNum] = LineFormat{
		typ:      "Rep",
		delim:    delim,
		values:   values,
		template: s.lines[lineNum-1],
	}
	return s
}

// Gen calls generator repeatedly until either
// [LineFormatter] notEmpty return value is false or the
// max iterations is reached. If notEmpty is false on the
// first call, the line is removed completely.
func (s *Sprintl) Gen(
	lineNum int,
	delim string,
	max int,
	generator LineFormatter,
) *Sprintl {
	s.formatters[lineNum] = LineFormat{
		typ:       "Gen",
		delim:     delim,
		max:       max,
		generator: generator,
		template:  s.lines[lineNum-1],
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
