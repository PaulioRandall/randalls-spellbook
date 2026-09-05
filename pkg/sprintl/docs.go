// Package sprintl is a simple multiline string formatter,
// hence the name. It was created to make SQL query string
// formatting easier and more dev friendly, but the package
// can be used for any kind of shortish multiline strings.
// It is intended for quick locally scoped method chaining,
// but the [Sprintl] core type can be passed around if
// needed.
//
//	"There are no solutions, only trade-offs."
//	                           – Thomas Sowell
//
// Existing SQL query builders force you to build the whole
// query via functions. While this greatly reduces the
// potential syntax errors and allows for lots of
// useful features, it also adds layers of abstraction that
// make debugging logical errors more difficult. Sprintl
// was designed so I could take a formatter's approach to
// creating SQL strings, rather than a builder's approach.
// I find being able to see the template and formatting it
// with a handful of simple functions more intuitive and
// easier to read, write, and debug.
//
// Primary formatting functions:
//   - [Sprintl.Fmt]
//   - [Sprintl.Dup]
//   - [Sprintl.Rep]
//   - [Sprintl.Gen]
//
// Supportive formatting functions (used straight after a
// primary function):
//   - [Sprintl.Join]
//   - [Sprintl.Marry]
//   - [Sprintl.Trim]
//
// General formatting functions applied to the formatted
// string of lines:
//   - [Sprintl.TrimSpace]
//   - [Sprintl.TrimLines]
//   - [Sprintl.PruneLines]
//
// Run formatting and get the result string:
//   - [Sprintl.String]
package sprintl
