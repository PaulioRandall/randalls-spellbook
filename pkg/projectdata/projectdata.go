// Package projectdata defines the projects core data
// structures.
//
// It does not validate data content, except maybe some ID
// integrity checking. Data validation is done in the
// datavalidation package (who saw that coming).
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
// flexible, and modular. It also simplifies use with the
// SQLite (Relational Database) storage package.
package projectdata

// ProjectData stores and provides access to project data.
type ProjectData struct {
	videos            []video
	videoStates       []videoState
	videoSections     []videoSection
	videoObservations []videoObservation
}
