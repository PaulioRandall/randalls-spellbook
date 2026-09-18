package data2

import (
	"github.com/google/uuid"
)

// randomEntityId randomly generates a new entity ID string
// in the form of a UUID.
func randomEntityId() string {
	return uuid.New().String()
}
