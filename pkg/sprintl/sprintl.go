package sprintl

import (
	"strings"
)

// LineFormatter is passed to the [Sprintl.G] function to
// enable customised formatting.
type LineFormatter func(
	formatter Formatter,
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
	formatters map[int]Formatter
}

// Lines returns a new [Sprintl] object for formatting the
// passed lines.
func Lines(lines ...string) *Sprintl {
	return &Sprintl{
		lines:      lines,
		formatters: map[int]Formatter{},
	}
}

// F formatters the specified line using the passed values.
func (s *Sprintl) F(
	lineNum int,
	values ...any,
) *Sprintl {
	s.formatters[lineNum] = Formatter{
		typ:      "F",
		values:   values,
		template: s.lines[lineNum-1],
	}
	return s
}

// R repeats the line for each value in values then appends
// the delim to each, except the last line. The line
// must contain a single placeholder, e.g. '%v'. Passing
// an empty values array removes the line completely.
func (s *Sprintl) R(
	lineNum int,
	delim string,
	values ...any,
) *Sprintl {
	s.formatters[lineNum] = Formatter{
		typ:      "R",
		delim:    delim,
		values:   values,
		template: s.lines[lineNum-1],
	}
	return s
}

// G calls generator repeatedly until either
// [LineFormatter] notEmpty return value is false or the
// max iterations is reached. If notEmpty is false on the
// first call, the line is removed completely.
func (s *Sprintl) G(
	lineNum int,
	delim string,
	max int,
	generator LineFormatter,
) *Sprintl {
	s.formatters[lineNum] = Formatter{
		typ:       "G",
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
