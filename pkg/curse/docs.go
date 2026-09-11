// Package curse is an alternative to the standard 'errors'
// package.
//
// Curse was created to generalise common error creation
// functions and provide better more readable terminal
// printing.
//
// TODO: Ensure interface is super set of errors so users
//
//	can swap out 'errors' for 'curse' and vice versa
//	easily.
//
// TODO: Use runtime.Caller(1) to get line number and
//
//	file.
package curse
