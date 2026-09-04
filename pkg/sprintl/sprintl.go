package sprintl

import (
	"strings"
)

// LineGenerator is used to format lines via the G
// function.
type LineGenerator func(
	index int,
	templateLine string,
) (
	result string,
	isEmpty bool,
)

// Sprintl is the formatter returned by Lines. It's public
// so you can pass it around if need, but it's intended
// for quick local use through method chaining.
type Sprintl struct {
	lines      []string
	formatters map[int]lineFormatter
}

// Lines returns a new Sprintl instance for formatting
// the passed lines.
//
//	var spl Sprintl := Lines(
//		"SELECT",
//		"  %s", // Line 2
//		"FROM",
//		"  %s", // Line 4
//	)
func Lines(lines ...string) *Sprintl {
	return &Sprintl{
		lines:      lines,
		formatters: map[int]lineFormatter{},
	}
}

// Compile compiles each line, joins them together, and
// returns the resultant string.
func (spl *Sprintl) Compile() string {
	lines := []string{}

	for i, line := range spl.lines {
		lineNum := i + 1

		form, ok := spl.formatters[lineNum]
		if !ok {
			lines = append(lines, line)
			continue
		}

		nextLines := form.apply(line)
		lines = append(lines, nextLines...)
	}

	return strings.Join(lines, "\n")
}

// String is an alias for Compile.
func (spl *Sprintl) String() string {
	return spl.Compile()
}

// F formatters the specified line using the passed args.
//
//	sql := Lines(
//		"SELECT",
//		"  %s || '-' || %s AS date,", // Line 2
//		"  revenue"
//		"FROM",
//		"  data",
//	).
//		F(2, "year", "month").
//		String()
//
//	sql == `SELECT
//	  year-month AS date,
//	  revenue
//	FROM
//		data`
func (spl *Sprintl) F(
	lineNum int,
	args ...any,
) *Sprintl {
	spl.formatters[lineNum] = lineFormatter{
		typ:  "F",
		args: args,
	}
	return spl
}

// R repeats the line for each value in values. The line
// must contain a single placeholder, e.g. '%v'. Passing
// an empty values array removes the line completely.
//
//	s := Lines(
//		"  %s", // Line 1
//	).
//		R(1, "line 1", "line 2", "line 3").
//		String()
//
//	s == `line 1
//	line 2
//	line 3`
func (spl *Sprintl) R(
	lineNum int,
	values ...any,
) *Sprintl {
	spl.formatters[lineNum] = lineFormatter{
		typ:    "R",
		values: values,
	}
	return spl
}

// RF repeats the line for each argument set in valueArgs.
// Passing an empty valueArgs array removes the line
// completely.
func (spl *Sprintl) RF(
	lineNum int,
	valueArgs ...[]any,
) *Sprintl {
	spl.formatters[lineNum] = lineFormatter{
		typ:       "RF",
		valueArgs: valueArgs,
	}
	return spl
}

// J repeats the line for each value in values then appends
// the delim to each, except the last line. The line
// must contain a single placeholder, e.g. '%v'. Passing
// an empty values array removes the line completely.
//
//	sql := Lines(
//		"SELECT",
//		"  %s", // Line 2
//		"FROM",
//		"  users",
//	).
//		J(2, ",", "name", "age", "height").
//		String()
//
//	sql == `SELECT
//	  name,
//	  age,
//		height
//	FROM
//		users`
func (spl *Sprintl) J(
	lineNum int,
	delim string,
	values ...any,
) *Sprintl {
	spl.formatters[lineNum] = lineFormatter{
		typ:    "J",
		delim:  delim,
		values: values,
	}
	return spl
}

// JF repeats the line for each argument set in valueArgs
// then appends the delim to each, except the last line.
// Passing an empty valueArgs array removes the line
// completely.
func (spl *Sprintl) JF(
	lineNum int,
	delim string,
	valueArgs ...[]any,
) *Sprintl {
	spl.formatters[lineNum] = lineFormatter{
		typ:       "JF",
		delim:     delim,
		valueArgs: valueArgs,
	}
	return spl
}

// G calls generator repeatedly until isEmpty is false.
// For each result, where isEmpty is true, the current
// index count and template line are passed. If isEmpty
// is false on the first call, the line is removed
// completely.
//
// Ensure your LineGenerator function is well tested as
// there is a notable risk of creating an infinite loop
// with this approach. Use the GN function where possible.
//
//	genYears := func(i int, line string) (string, bool) {
//		if i >= 5 {
//			return "", false
//		}
//
//		year := fmt.Sprintf(line, 2020 + i)
//		if i < 4 {
//			year += ","
//		}
//
//		return year, true
//	}
//
//	sql := Lines(
//		"SELECT",
//		"  %s", // Line 2
//		"FROM",
//		"  sales_data",
//	).
//		G(2, genYears).
//		String()
//
//	sql == `SELECT
//	  2020,
//	  2021,
//		2022,
//		2023,
//		2024
//	FROM
//		sales_data`
func (spl *Sprintl) G(
	lineNum int,
	generator LineGenerator,
) *Sprintl {
	spl.formatters[lineNum] = lineFormatter{
		typ:       "G",
		generator: generator,
	}
	return spl
}

// GN calls generator repeatedly n times. A result is
// ignored if isEmpty is true, the current
// index count and template line are passed. If the
// generator's isEmpty return value is false, but the
// generation will continue. If n is less than 1 the line
// is removed completely.
//
//	genYears := func(i int, line string) (string, bool) {
//		return fmt.Sprintf(line, 2020 + i), true
//	}
//
//	sql := Lines(
//		"SELECT",
//		"  %s,", // Line 2
//		"  average",
//		"FROM",
//		"  sales_data",
//	).
//		G(2, 5, genYears).
//		String()
//
//	sql == `SELECT
//	  2020,
//	  2021,
//		2022,
//		2023,
//		2024,
//		average
//	FROM
//		sales_data`
func (spl *Sprintl) GN(
	lineNum int,
	n int,
	generator LineGenerator,
) *Sprintl {
	spl.formatters[lineNum] = lineFormatter{
		typ:            "GN",
		generator:      generator,
		generatorCount: n,
	}
	return spl
}

/*
// TODO
//
// JG calls generator repeatedly until either notEmpty is
// false or the max iterations is reached. If notEmpty is
// false on the first call, the line is removed completely.
// Like the J function, the delim is appended to each line
// except the last.
//
// If your use case requires a few thousand or more lines
// then I don't recommend using this package.
func (spl *Sprintl) JG(
	lineNum int,
	delim string,
	max int,
	generator LineGenerator,
) *Sprintl {
	spl.formatters[lineNum] = lineFormatter{
		typ:            "JG",
		delim: delim,
		max: max,
		generator:      generator,
	}
	return spl
}
*/
