package entity

import (
	"github.com/google/uuid"
)

type Team struct {
	ID   uuid.UUID
	Name string
}
