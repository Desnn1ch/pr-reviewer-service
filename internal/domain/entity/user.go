package entity

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	TeamID   uuid.UUID
	Name     string
	IsActive bool
}
