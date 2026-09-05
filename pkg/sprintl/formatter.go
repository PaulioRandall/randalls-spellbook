package sprintl

import (
	"fmt"
)

// Formatter holds information about a line formatting.
// It is passed to [LineFormatter] functions when using the
// [Sprintl.G] function.
type Formatter struct {
	typ       string
	delim     string
	values    []any
	max       int
	generator LineFormatter
	index     int
	template  string
}

// Index returns the current index in the [Sprintl.G] loop.
func (f Formatter) Index() int {
	return f.index
}

// Delim returns the delimiter, which may be an empty
// string.
func (f Formatter) Delim() string {
	return f.delim
}

// Max returns the max number of times a [LineFormatter]
// function will be called for a particular line. Only
// applicable when using the [Sprintl.G] function.
func (f Formatter) Max() int {
	return f.max
}

// Template returns the template line provided to the
// [Lines] function.
func (f Formatter) Template() string {
	return f.template
}

// Fmt is a convenience function for calling fmt.Sprintf
// with the template as the template string.
func (f Formatter) Fmt(args ...any) string {
	return fmt.Sprintf(f.template, args...)
}

func (f Formatter) apply() []string {
	switch f.typ {
	case "F":
		return f.applyF()
	case "R":
		return f.applyR()
	case "G":
		return f.applyG()
	default:
		msg := fmt.Sprintf("Unknown format type '%s'", f.typ)
		panic(msg)
	}
}

func (f Formatter) applyF() []string {
	return []string{
		fmt.Sprintf(f.template, f.values...),
	}
}

func (f Formatter) applyR() []string {
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

func (f Formatter) applyG() []string {
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

func (f Formatter) format(
	template string,
	value any,
) string {
	if args, ok := value.([]any); ok {
		return fmt.Sprintf(template, args...)
	} else {
		return fmt.Sprintf(template, value)
	}
}
