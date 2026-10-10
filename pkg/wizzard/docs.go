// Package wizzard provides a simple ORM-based interface
// for managing SQLite databases with Go struct's.
//
// # Rincewind the Wizzard
//
// No, wizzard is not a mis-spelling of wizard! Well ok
// it is, but it's intentional. Wizzard is an old English
// rock band that are famous for the song "I Wish It Could
// be Christmas Everyday". But this package is actually
// named after the words on Rincewind's hat which was
// created using sequins. Rincewind is a failed wizard of
// the Discworld who lives in a never-ending series of
// interesting times. They say that if you've never read a
// Discworld novel you're not really real, but since very
// few people are really real, 'they' included, I wouldn't
// worry about it. But I do recommend reading the entire
// Discworld series.
//
// # Models
//
// Models are created from an object's type (Go type)
// information and represent a mapping between a Go type
// and a SQLite database table. Go types and SQLite tables
// are completely decoupled. A model may be partial.
// Partial models don't map field-to-column perfectly with
// some exported struct fields that don't appear in the
// table, and table columns that don't appear as fields in
// the struct.
//
// This package does not create or manage database
// connections, it only uses them. Connections must be
// passed to various functions to create models and operate
// on the database.
//
// The SQL functions that operate on the database (e.g.
// [Create], [Insert], [Select], etc) are intentionally
// standalone functions with no receiver. This allows them
// to be copied and modified easily. There's also functions
// utility functions for scanning rows into the model's
// type (i.e. [ScanRows] and [ScanFirstRow]).
package wizzard
