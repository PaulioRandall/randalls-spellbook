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
// and [Storm.Drop] us the type information to determine
// what to create or drop. The object and table are
// completely decoupled, that is, any object with the
// same name can be used as a model or partial model for
// inserting and selecting after the table is created.
//
// # Partial models
//
// Partial models have the same name as a table but done
// have all the same exported fields with compatible field
// types. For some operations such as inserting, updating,
// and selecting they can be used to provide or access a
// subset of columns from the table. It's a similar concept
// to Go's interfaces, as long as the model's name matches
// it counts as a model for that table.
//
// # Objects
//
// Objects are passed for their type and data. Inserting
// and updating only use the type information of the value
// passed them to determine what tables and columns to
// insert or update.
//
// # Creating tables
//
// Tables are explicitly created using the [Storm.Create]
// function and may be implicitly created when using the
// any insert functions. The name of the model (struct
// type) is the name of the table, and the exported field
// names and types determine the column names and types
// repectively. The first field is designated the PRIMARY
// KEY regardless of name and type. All columns are
// NOT NULL and the DEFAULT value will be the zero value of
// field's type. However, only primitive types may be used
// as field types:
//
//	INTEGER:
//		int, int8, int16, int32, int64,
//		uint, uint8, uint16, uint32, uint64
//	REAL:
//		float32, float64
//	TEXT:
//		string
//
// Dropping tables
//
//	TODO
//
// # Inserting data
//
// Insert data by passing the objects you want stored to
// [Storm.Insert]. If the table doesn't exist then the
// object's type (model) information will be used to create
// it first. You can also pass objects of partial models to
// insert a subset of data; the other columns will default
// to their zero values.
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
