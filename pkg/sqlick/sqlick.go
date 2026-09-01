package sqlick

import (
	"errors"
)

var (
	// ErrNoTableForType is used when an object is passed
	// who's type does not map to a table. This means the
	// type was not added via Sqlick.AddStructTable before
	// use.
	ErrNoTableForType = errors.New(
		"No matching table for object type",
	)
)

type Sqlick interface {
	Open() error
	IsOpen() bool
	Close() error
	AddStructTable(any) error
	CreateTables() error
	Insert(any) error
}
