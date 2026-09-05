package sprintl

import (
	"fmt"
	"strings"
)

// LineGenerator is passed to the [Sprintl.Gen] to enable
// bespoke formatting.
type LineGenerator func(
	formatter LineFormatter,
) (
	result string,
	notEmpty bool,
)

// LineFormatter holds information about a line's
// formatting needs. It is passed to [LineGenerator]
// functions when using the [Sprintl.Gen] function.
type LineFormatter struct {
	typ         string
	prefix      string
	delim       string
	prefixDelim bool
	trim        bool
	values      []any
	max         int
	generator   LineGenerator
	index       int
	template    string
}

// Index returns the current repetition index.
func (lf LineFormatter) Index() int {
	return lf.index
}

// Prefix returns the prefix that should be prepended to
// every line in a repetition. It may be empty.
func (lf LineFormatter) Prefix() string {
	return lf.prefix
}

// Delim returns the join delimiter. It may be empty.
func (lf LineFormatter) Delim() string {
	return lf.delim
}

// PrefixDelim returns true if the value returned by
// [LineFormatter.Delim] should be applied to the front of
// all lines except the first, instead of to the end of all
// lines except the last.
func (lf LineFormatter) PrefixDelim() bool {
	return lf.prefixDelim
}

// Max returns the max number of times a [LineGenerator]
// function will be called for a particular line.
func (lf LineFormatter) Max() int {
	return lf.max
}

// Template returns the template line of the current
// iteration. It will be line indexed by
// [LineFormatter.Index] in the argument passed to [Lines]
// or [Split].
func (lf LineFormatter) Template() string {
	return lf.template
}

// Fmt is a convenience function for calling fmt.Sprintf
// with the template as the template string i.e. only pass
// the args to this function.
func (lf LineFormatter) Fmt(args ...any) string {
	return fmt.Sprintf(lf.template, args...)
}

func (lf LineFormatter) apply() []string {
	switch lf.typ {
	case "Fmt":
		return lf.applyFmt()
	case "Dup":
		return lf.applyDup()
	case "Rep":
		return lf.applyRep()
	case "Gen":
		return lf.applyGen()
	default:
		msg := fmt.Sprintf("Unknown format type '%s'", lf.typ)
		panic(msg)
	}
}

func (lf LineFormatter) applyFmt() []string {
	lines := []string{
		fmt.Sprintf(lf.template, lf.values...),
	}
	lf.trimLine(lines)
	return lines
}

func (lf LineFormatter) applyDup() []string {
	lines := make([]string, lf.max, lf.max)

	for i := 0; i < lf.max; i++ {
		lf.index = i
		lines[i] = lf.format(lf.template, lf.values)
		lf.applyLineMods(lines)
		lf.trimLine(lines)
	}

	return lines
}

func (lf LineFormatter) applyRep() []string {
	lineCount := len(lf.values)
	lines := make([]string, lineCount, lineCount)

	for i, v := range lf.values {
		lf.index = i
		lines[i] = lf.format(lf.template, v)
		lf.applyLineMods(lines)
		lf.trimLine(lines)
	}

	return lines
}

func (lf LineFormatter) applyGen() []string {
	lines := []string{}

	for i := 0; i < lf.max; i++ {
		lf.index = i
		line, ok := lf.generator(lf)

		if !ok {
			break
		}

		lines = append(lines, line)
		lf.applyLineMods(lines)
		lf.trimLine(lines)
	}

	return lines
}

func (lf LineFormatter) format(
	template string,
	value any,
) string {
	if args, ok := value.([]any); ok {
		return fmt.Sprintf(template, args...)
	} else {
		return fmt.Sprintf(template, value)
	}
}

func (lf LineFormatter) applyLineMods(lines []string) {
	i := lf.index

	if i <= 0 {
		lines[i] = lf.prefix + lines[i]
		return
	}

	if lf.prefixDelim {
		lines[i] = lf.prefix + lf.delim + lines[i]
		return
	}

	lines[i-1] += lf.delim
	lines[i] = lf.prefix + lines[i]
}

func (lf LineFormatter) trimLine(lines []string) {
	if lf.trim {
		i := lf.index
		lines[i] = strings.TrimSpace(lines[i])
	}
}
