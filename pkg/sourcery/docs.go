// Package sourcery is a WebView based desktop application
// framework focusing on simple desktop applications and
// rapid prototyping.
//
// # TODO
//
//   - Now that I have a copy of go-webview it may be
//     worth implanting the Thonk package there. Either
//     add a new function to the WebView interface,
//     make functions passed to bind accessible through
//     JS calls to the Go function, or change
//     the Bind function users must make all calls
//     through the Go function. I think the later is
//     preferable because it simplifies the interface.
//     Yes it reduces freedom of expression but I think
//     having all backend user calls be directed
//     through a single JS function will make for
//     clearer frontends.
//   - Implanting Sourcery package there also bit by bit.
package sourcery
