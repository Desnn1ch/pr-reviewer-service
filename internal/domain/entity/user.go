package entity

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	TeamName string
	Name     string
	IsActive bool
}
