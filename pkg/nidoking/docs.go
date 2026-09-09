// Package nidoking provides string searching and
// replacement with additional focus on line operations.
//
// The lower level API provides string searching and
// replacement, i.e. does what the standard 'strings'
// package does but tailored for common use cases and
// provides information rich results with minimal code. The
// main drawback is that multi-line search terms (needles)
// are currently not supported. The primary functions and
// types are [Find], [NeedleInHaystack], and
// [Replacement].
//
// The higher level API provides simple string template
// formatting. It's much simpler than the standard
// 'text/template' but package also less feature rich.
// Dynamic values are defined using mustache tokens,
// e.g. '{{key}}', the same as 'text/template' so some code
// editors will highlight tokens within the strings.
// However, they are just placeholders; how values are
// inserted depends on the template functions you call
// and the arguments you give it, and Some functions
// require pairs of tokens. A [Template] can be created
// via the [Given] and [Lines] functions.
//
// Both APIs were designed for locally scoped usage with
// method chaining but types are available for passing
// around or creating adapters.
package nidoking
