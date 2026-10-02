// Package storm is a minimalistic ORM based database
// API using type information of struct's.
//
// Instead of writing or building query statements, user
// programmers pass structs to represent tables or part of
// a table. Objects with these struct types are passed to
// various functions to perform standard database
// actions. It's an approach used by existing database
// packages, such as https://gorm.io/.
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
// Models are objects passed for their type information.
// Their actual data is ignored. Some operations only use
// the type information of the passed object while others
// use the type information and the object data. Structs
// and tables are completely decoupled, i.e. any struct
// may be used as a model or partial model for operations
// such as inserting, selecting, and deleting from a table;
// the name of the struct determines which table they
// operate on.
//
// Partial models will have the same name as a table but
// not all the same exported fields. For some operations
// such as inserting, updating, and selecting they can be
// used to provide or access a subset of columns from the
// table. It's a similar concept to Go's interfaces, as
// long as the model's name matches then it counts as a
// model for that table.
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
package storm
