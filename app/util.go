package app

import (
	"github.com/google/uuid"
)

func randomEntityId() string {
	return uuid.New().String()
}
