package sprintl

import (
	"fmt"
)

// LineFormat holds information about a line formatting.
// It is passed to [LineFormatter] functions when using the
// [Sprintl.Gen] function.
type LineFormat struct {
	typ       string
	delim     string
	values    []any
	max       int
	generator LineFormatter
	index     int
	template  string
}

// Index returns the current index in the [Sprintl.Gen]
// loop.
func (f LineFormat) Index() int {
	return f.index
}

// Delim returns the delimiter, which may be an empty
// string.
func (f LineFormat) Delim() string {
	return f.delim
}

// Max returns the max number of times a [LineFormatter]
// function will be called for a particular line. Only
// applicable when using the [Sprintl.Gen] function.
func (f LineFormat) Max() int {
	return f.max
}

// Template returns the template line provided to the
// [Lines] function.
func (f LineFormat) Template() string {
	return f.template
}

// Fmt is a convenience function for calling fmt.Sprintf
// with the template as the template string.
func (f LineFormat) Fmt(args ...any) string {
	return fmt.Sprintf(f.template, args...)
}

func (f LineFormat) apply() []string {
	switch f.typ {
	case "Fmt":
		return f.applyF()
	case "Rep":
		return f.applyR()
	case "Gen":
		return f.applyG()
	default:
		msg := fmt.Sprintf("Unknown format type '%s'", f.typ)
		panic(msg)
	}
}

func (f LineFormat) applyF() []string {
	return []string{
		fmt.Sprintf(f.template, f.values...),
	}
}

func (f LineFormat) applyR() []string {
	lineCount := len(f.values)
	lines := make([]string, lineCount, lineCount)

	for i, v := range f.values {
		if i > 0 {
			lines[i-1] += f.delim
		}

		lines[i] = f.format(f.template, v)
	}

	return lines
}

func (f LineFormat) applyG() []string {
	lines := []string{}

	for i := 0; i < f.max; i++ {
		f.index = i
		line, ok := f.generator(f)

		if !ok {
			break
		}

		if i > 0 {
			lines[i-1] += f.delim
		}

		lines = append(lines, line)
	}

	return lines
}

func (f LineFormat) format(
	template string,
	value any,
) string {
	if args, ok := value.([]any); ok {
		return fmt.Sprintf(template, args...)
	} else {
		return fmt.Sprintf(template, value)
	}
}
