package app

import (
	"context"

	"github.com/google/uuid"

	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/entity"
)

type TxManager interface {
	InTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type TeamRepo interface {
	Create(ctx context.Context, team entity.Team) error
	GetByName(ctx context.Context, name string) (entity.Team, error)
}

type UserRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (entity.User, error)
	ListByTeamName(ctx context.Context, teamName string) ([]entity.User, error)
	ListActiveByTeamName(ctx context.Context, teamName string) ([]entity.User, error)
	UpsertMany(ctx context.Context, users []entity.User) error
	SetActive(ctx context.Context, userID uuid.UUID, isActive bool) error
}

type PRRepo interface {
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	Create(ctx context.Context, pr entity.PR) error
	GetByID(ctx context.Context, id uuid.UUID) (entity.PR, error)
	Update(ctx context.Context, pr entity.PR) error
	ListByReviewerID(ctx context.Context, reviewerID uuid.UUID) ([]entity.PR, error)
}
