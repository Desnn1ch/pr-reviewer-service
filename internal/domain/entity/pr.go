package entity

import (
	"time"

	"github.com/google/uuid"
)

type PRStatus string

const (
	StatusOpen   PRStatus = "OPEN"
	StatusMerged PRStatus = "MERGED"
)

type PR struct {
	ID        uuid.UUID
	Title     string
	AuthorID  uuid.UUID
	Status    PRStatus
	CreatedAt time.Time
	MergedAt  *time.Time

	Reviewers []uuid.UUID
}

func (p PR) CanChangeReviewers() bool { return p.Status == StatusOpen }
func (p PR) IsMerged() bool           { return p.Status == StatusMerged }
