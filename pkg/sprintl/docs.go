// Package sprintl is a multiline string formatter, hence
// the name. It was created to make SQL query string
// formatting easier and more dev friendly, but the package
// can be used for any kind of shortish multiline strings.
//
// It is intended for quick locally scoped method chaining,
// but the [Sprintl] core type is available if you need to
// pass it around.
//
// Existing SQL query builders force you to build the whole
// query via functions. While this greatly reduces the
// probability of syntax errors and allows for lots of
// features, it also adds layers of abstraction that make
// debugging logical errors more difficult. Sprintl was
// designed so I could take a formatters approach to
// creating SQL strings, rather than a builders approach.
// I find the visible template lines more intuitive and
// easier to read, write, and debug with.
package sprintl
