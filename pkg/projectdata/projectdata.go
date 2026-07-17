// Package projectdata manages a project's data and data
// storage.
//
// Data structures use a Referential Modeling approach to
// dependency and hierarchy. This is the same approach used
// by relational databases. The structure is flat with
// every entity having a unique ID (primary key) which is
// referenced by dependent entities (foreign key).
// Functions are intentionally kept simple with individual
// CRUD operations for each entity type.
//
// I (Paulio) choose this approach because it is simple,
// flexible, and modular so I can . It also simplifies use
// with the data storage layer that uses SQLite (Relational
// Database).
//
// The data layer is the source of truth for all project
// data. The data layer also provides a subscription
// service allowing notification of data changes.
package projectdata

// ProjectData stores and provides access to project data.
type ProjectData struct {
	videos            []Video
	videoStates       []VideoState
	videoSections     []VideoSection
	videoObservations []VideoObservation
}
