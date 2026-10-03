// Package storm is a minimalistic ORM based SQLite
// database API using type information of struct's.
//
// Instead of writing or building query statements, user
// programmers pass objects (instances of structs) to
// represent tables or part of a table. Objects with these
// struct types are passed to various functions to perform
// database operations. It's an approach used by existing
// database packages such as https://gorm.io/.
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
// of any table; whether it is suitable is your descision.
// A model may be partial. Partial models don't map
// field-to-column perfectly with some exported struct
// fields that don't appear in the table, and table columns
// that don't appear as a field in the model's Go type.
//
// # Operations
//
//   - [Storm.Create] explicitly creates tables.
//   - [Storm.Put] inserts or updates (upserts) data. If
//     the table doesn't exist then the object's type
//     (model) information will be used to create it. You
//     can also pass objects of partial models to insert or
//     update a subset of data; on insert the other columns
//     will default to their zero values.
//   - [Storm.List] to select all data for a specific table.
//   - [Storm.Get] to select a specific entry by ID.
//   - [Storm.Delete] to delete an entry by ID.
//   - [Storm.Drop] to remove a table. All data is deleted
//     in the process and there's no way to restore it. To
//     protect data, create regular backups of the database
//     file.
//
// # TODO
//   - AS operations, e.g. CreateAs("Users", Model{})
//     Allow table name to be specified rather than using
//     the struct's name.
//   - Function to perform custom operations. Must lock
//     and allow access to cachedMapper and db.
//   - Be more targetted with mutex use. Only lock when
//     performing Model parsing and caching. Ponder use
//     of RWMutex instead.
//   - Create operation modes (could create adapters using
//     an interface instead so structs representing the
//     below modes are created instead of using modes):
//   - OperationModeFree: tables are created automatically
//     for Put operations, List and Get operations return
//     an empty/zero result set or result object if the
//     table doesn't exist (not an error), Create
//     returns without error if the table already exists,
//     and Drop and Delete operations return without error
//     if the table doesn't exist.
//   - OperationModeStrict: tables are not created
//     automatically for Put operations (an error is
//     returned instead), List and Get operations return
//     an error if the table or speecific entry doesn't
//     exist, Create returns an error if the table already
//     exists, and Drop and Delete operations return an
//     error if the table doesn't exist.
//   - Create function to copy or backup database.
//   - Rename package to stormy
package storm
