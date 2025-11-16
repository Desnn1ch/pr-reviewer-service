package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	TeamID    uuid.UUID
	Name      string
	IsActive  bool
	CreatedAt time.Time
}
