// Package data manages a project's data and data storage.
//
// Data structures use a Referential Modeling approach to
// dependency and hierarchy. This is the same approach used
// by relational databases. The structure is flat with
// every entity having a unique ID (primary key) which is
// referenced by dependent entities (foreign key).
//
// I (Paulio) choose this approach because it is flexible
// and modular allowing transformations into other formats
// with relative ease. It also simplifies use with the data
// storage layer that uses SQLite (Relational Database).
//
// The data layer is the source of truth for all project
// data. As a general rule, entity IDs should be stored,
// passed, and used within the application, then data
// requests made when data is needed. This should minimise
// issues with stale data. The data layer also provides a
// subscription service allowing notification of data
// changes.
package data

// EntityID is the core ID type that every enitity, within
// the project, will have a unique value for.
type EntityID string
