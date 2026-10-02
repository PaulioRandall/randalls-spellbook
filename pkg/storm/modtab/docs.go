// Package modtab provides a simple ORM-based interface to
// managing SQLite databases with Go struct's.
//
// # TODO
//   - Consider renaming package to 'wizzard'.
//   - Update and tidy storm package before applying ideas
//     below this one.
//   - Allow table name to differ from struct name by
//     creating ParseAs and MapAs functions.
//   - Add Insert and Update functions to Model.
//   - Add Select function to Model that accepts string
//     WHERE clause followed by series of arguments to
//     SQL parameters and returns all results.
//   - Add SelectFirst function to Model that accepts
//     string  WHERE clause followed by series of arguments to
//     SQL parameters and returns a single result.
package modtab
