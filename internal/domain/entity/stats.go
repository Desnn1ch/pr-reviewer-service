package entity

import "github.com/google/uuid"

type ReviewerStats struct {
	UserID          uuid.UUID
	Username        string
	TeamName        string
	AssignedOpenPRs int
}
