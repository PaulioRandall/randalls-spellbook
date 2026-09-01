package sqlick

import (
	"errors"
)

var (
	// ErrNoTableForType is used when an object is passed
	// to a Sqlick function who's type does not map to a
	// table. This means the type was not added via
	// Sqlick.Register before use.
	ErrNoTableForType = errors.New(
		"No matching table for object type",
	)
)

type Sqlick interface {
	Open() error
	IsOpen() bool
	Close() error
	Register(any) error
	CreateTables() error
	Insert(any) error
	Update(any) error
}
