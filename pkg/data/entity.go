package data

import (
	"github.com/google/uuid"
)

// randomEntityId randomly generates a new entity ID string
// form of a UUID.
func randomEntityId() string {
	return uuid.New().String()
}
