// Package stormy is a minimalistic ORM based SQLite
// database API.
//
// Instead of writing or building query statements, user
// programmers pass objects to manage table structure and
// content. It's an approach used by existing database
// packages such as https://gorm.io/.
//
// This package provides a minimalistic solution that
// trades away the feature richness and configurability of
// existing tools for speed and ease of use when developing
// simple desktop apps and rapid prototypes. It also aims
// for high plunderability, i.e. the ability for people to
// copy and modify the codebase for their own purposes.
//
// # Models
//
// Models are created from an object's type (Go type)
// information and represent a mapping between a Go type
// and a SQLite database table. Go types and SQLite tables
// are completely decoupled, i.e. a model may be
// constructed from a Go type which is then used for
// database operations. A single Go type may be used for
// multiple tables using the functions ending in 'As', e.g.
// CreateAs. A model may be partial. Partial models don't
// map field-to-column perfectly with some exported struct
// fields that don't appear in the table, and table columns
// that don't appear as fields in the struct.
//
// # Operations
//
// Operations were designed to be minimise errors by
// making positive assumptions. For example, a check is
// made on the existance of a table before executing a
// query. If the table doesn't exist then an appropriate
// action is taken, for instance, [Stormy.Put] will create
// the table so it has somewhere to insert data, while
// [Stormy.Get] will just return false to indicate the
// item was not found rather than returning an error.
//
//   - [Stormy.Create] and [Stormy.CreateAs] explicitly
//     create tables.
//   - [Stormy.Put] and [Stormy.PutAs] insert or update
//     (upsert) data.
//   - [Stormy.List] and [Stormy.ListAs] to select all
//     objects matching criteria.
//   - [Stormy.Get] and [Stormy.GetAs] to select the first
//     object matching criteria.
//   - [Stormy.Delete] and [Stormy.DeleteAs] to delete all
//     objects matching criteria.
//   - [Stormy.Drop] and [Stormy.DropAs] to remove a table.
//   - [Stormy.Custom] allows custom queries.
//
// # TODO
//   - Add new mode: OpenMode,
//   - OpenModeManual: that means the database performs
//     no automated opening or closing of the database.
//   - OpenModeRequest: upon making an API call that
//     requires an open DB, if it is already open then
//     it will remain open after the request, if it is
//     not open then it will be opened for the request then
//     closed.
//   - OpenModePersist: upon making an API call that
//     requires an open DB, if it is closed then it will be
//     opened and remain open after the request (DEFAULT)
package stormy
