package stormy

import (
	"github.com/PaulioRandall/randalls-spellbook/pkg/sin"
)

var (
	// ErrForDatabase is returned for almost all errors and
	// prints the database path.
	ErrForDatabase = sin.Template(
		"Stormyy database error '%s'",
	)

	// ErrForTable occurs in the chain of every error
	// produced from an operation on a known table.
	ErrForTable = sin.Template(
		"For table '%s'",
	)

	// ErrForType occurs in the chain of every error
	// produced from an operation on a known model.
	ErrForType = sin.Template(
		"For model '%s'",
	)

	// ErrNotOpen occurs when trying to perform an operation
	// before opening the database.
	ErrNotOpen = sin.Err(
		"Database not open",
	)

	// ErrRowScan is returned when an error occurs
	// scanning database results.
	ErrRowScan = sin.Template(
		"When scanning row '%d'",
	)

	// ErrNoIdField occurs when attempting to perform an
	// operation requiring a key (ID field) but the model
	// used does not contain it for the table being operated
	// on.
	ErrNoIdField = sin.Err(
		"Model missing key (ID field)",
	)
)
