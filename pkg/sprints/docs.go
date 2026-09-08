// Package sprints formats objects (instances of structs)
// into strings to aid debugging.
//
// The following rules apply to all public functions:
//   - Panic ensues when encountering any errors.
//   - Passing a value that doesn't derefernce to a struct is an error.
//   - Only exported fields are printed by default.
//   - Use OptionShowUnexported option to show unexported field names.
//   - When using OptionShowUnexported, the field value will always be <unexported>.
//   - Linefeeds in strings are printed in string form '\n'.
//   - Strings are sliced to a maximum of 30 runes + "...".
//   - Fields of nested structs are not printed.
//   - Zero value nested structs are printed as 'Name{}'.
//   - Non-zero nested structs are printed as 'Name{...}'.
//   - Non-empty arrays and slices are printed like '[3]Name{...}'.
//   - Empty arrays and slices are printed as '[0]Name{}'.
//   - Nil arrays and slices are printed as '[]Name'.
//   - Pointer values are dereferenced to show the actual value.
//   - Pointer values are preceeded with a single '*' for each level of indirection.
//   - Pointers to nil are printed as '⁎' (low asterisk) instead of '*'.
//   - Functions are currently printed as package fmt dictates.
//
// TODO: Print function definitions.
package sprints
