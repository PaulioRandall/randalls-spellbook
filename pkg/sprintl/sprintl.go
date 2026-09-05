package sprintl

import (
	"strings"
)

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
}

// Lines returns a new [Sprintl] object for formatting the
// passed lines.
func Lines(lines ...string) *Sprintl {
	return &Sprintl{
		lines:      lines,
		formatters: map[int]LineFormatter{},
	}
}

// Split returns a new [Sprintl] object for formatting the
// passed lines by splitting the passed string into lines.
func Split(s string) *Sprintl {
	return &Sprintl{
		lines:      strings.Split(s, "\n"),
		formatters: map[int]LineFormatter{},
	}
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
