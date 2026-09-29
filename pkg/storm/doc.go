// Package storm is a minimalistic ORM based database
// API for storing and accessing simple structures/tables.
//
// Instead of writing or building query statements, user
// programmers pass structs to represent tables or subsets
// of table rows. Objects with these struct types are
// passed to various functions to perform standard database
// actions. It's an approach used by existing database
// packages, such as https://gorm.io/.
//
// This package provides a minimalistic solution that
// trades away the feature richness and configurability of
// existing tools for ease of use when developing simple
// desktop tools and rapid prototypes. It also aims
// for high plunderability, i.e. the ability for people to
// copy and modify the codebase for their own purposes.
//
// # Models
//
// Models are objects passed for their type information.
// Their actual data is ignored. Calls to [Storm.Create]
// and [Storm.Drop] only use the type information of the
// object. Structs and tables are completely decoupled,
// that is, any struct with the same name may be used as a
// model or partial model for inserting and selecting from
// a table.
//
// Partial models may have the same name as a table but not
// all the same exported fields with compatible field
// types. For some operations such as inserting, updating,
// and selecting they can be used to provide or access a
// subset of columns from the table. It's a similar concept
// to Go's interfaces, as long as the model's name matches
// it counts as a model for that table.
//
// # Creating tables
//
// Tables are explicitly created using the [Storm.Create]
// function and may be implicitly created when using
// insert functions. See [Storm.Create] for more details.
//
// # Dropping tables
//
// Tables can be explicitly dropped using [Storm.Drop].
// All data is deleted in the process and there's no way to
// restore it. To protect data, create regular backups of
// the database file.
//
// # Inserting data
//
// Insert data by passing the objects you want stored to
// [Storm.Insert]. If the table doesn't exist then the
// object's type (model) information will be used to create
// it. You can also pass objects of partial models to
// insert a subset of data if the database already exists;
// the other columns will default to their zero values.
//
// # Updating data
//
//	TODO
//
// # Selecting data
//
//	TODO
//
// # Deleting data
//
//	TODO
package storm
