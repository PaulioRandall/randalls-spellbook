// Package storm is a minimalistic ORM based database
// API for storing and accessing simple structures/tables.
//
// Instead of writing or building query statements, users
// design structs with the storage structure they want.
// Objects with these struct types are passed to various
// functions to perform standard database actions such as
// creating tables and inserting data. It's an approach
// used by existing database packages, such as
// https://gorm.io/.
//
// This package provides a minimalistic solution that
// trades away the feature richness and configurability of
// existing tools for ease of use when developing simple
// desktop tools and rapid prototyping. It also aims
// for high plunderability, i.e. the ability for people to
// copy and modify the codebase for their own purposes.
//
// Throughout this package the term 'model' refers to an
// object passed purely for its type (create, select),
// while 'object' is used when the data is the subject
// (insert, update, delete).
package storm
