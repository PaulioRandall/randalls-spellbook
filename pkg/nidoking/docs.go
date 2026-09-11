// Package nidoking provides string searching and
// replacement with additional focus on line operations.
//
// The lower level API provides string searching and
// replacement, i.e. does what the standard 'strings' and
// 'regexp' packages do but tailored for common use cases,
// ease of use, and provides information rich results with
// minimal code. The main drawback is that multi-line
// search terms are not supported. The primary functions
// and types are [Find], [Match], [NeedleInHaystack], and
// [Replacement].
//
// The higher level API provides simple string template
// formatting. It's much simpler than the standard
// 'text/template' package but also much less feature rich.
// Dynamic values are defined using mustache tokens,
// e.g. '{{key}}' and {{key.fieldOrMethod}}, similar to
// 'text/template' so some code editors will highlight
// tokens within the strings. However, they are relatively
// simple placeholders; how values are inserted depends
// mostly on the template functions called and the
// arguments passed. A [Template] can be created via the
// [Given] and [Lines] functions. Template contains
// generic methods which unfortunately limits the
// interfacing capability.
//
// Both APIs were designed for locally scoped usage with
// method chaining but types are available for passing
// around and creating adapters.
package nidoking
