// Package storm is a minimalistic ORM based SQLite
// database API.
//
// Instead of writing or building query statements, user
// programmers pass objects (instances of structs) to
// represent tables or parts of them. Objects are passed
// to various functions to perform database operations.
// It's an approach used by existin database packages such
// as https://gorm.io/.
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
// Models are created from an object's type information and
// represent a mapping between a Go type and a SQLite
// database table. Go types and SQLite tables are
// completely decoupled, i.e. a struct may be used to
// create a model or partial model for database operations
// of any table; whether it is suitable is your decision.
// A single struct may be used to for multiple tables using
// the functions ending in 'As', e.g. CreateAs. A model
// may be partial. Partial models don't map field-to-column
// perfectly with some exported struct fields that don't
// appear in the table, and table columns that don't appear
// as a field in the model's Go type.
//
// # Operations
//
//   - [Storm.Create] and [Storm.CreateAs] explicitly
//     create tables.
//   - [Storm.Put] and [Storm.PutAs] insert or update
//     (upsert) data.
//   - [Storm.List] and [Storm.ListAs] to select all
//     objects matching criteria.
//   - [Storm.Get] and [Storm.GetAs] to select the first
//     object matching criteria.
//   - [Storm.Delete] and [Storm.DeleteAs] to delete all
//     objects matching criteria.
//   - [Storm.Drop] and [Storm.DropAs] to remove a table.
//   - [Storm.Custom] allows custom queries.
//
// # TODO
//   - Create function to copy or backup database.
//   - Rename package to stormy.
//   - Reorg named errors & test they are returned.
package storm
