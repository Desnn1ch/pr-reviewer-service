package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/Desnn1ch/pr-reviewer-service/internal/domain"
	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/entity"
)

type TeamRepo struct {
	db *DB
}

func NewTeamRepo(db *DB) *TeamRepo {
	return &TeamRepo{db: db}
}

func (r *TeamRepo) Create(ctx context.Context, team entity.Team) error {
	e := r.db.getExec(ctx)

	const q = `
		INSERT INTO teams (id, name, created_at)
		VALUES ($1, $2, $3)
`

	_, err := e.ExecContext(ctx, q,
		team.ID,
		team.Name,
		team.CreatedAt,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return domain.ErrTeamExists
		}
		return err
	}

	return nil
}

func (r *TeamRepo) GetByName(ctx context.Context, name string) (entity.Team, error) {
	e := r.db.getExec(ctx)

	const q = `
		SELECT id, name, created_at
		FROM teams
		WHERE name = $1
`

	var t entity.Team
	err := e.QueryRowContext(ctx, q, name).Scan(
		&t.ID,
		&t.Name,
		&t.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Team{}, domain.ErrNotFound
		}
		return entity.Team{}, err
	}

	return t, nil
}

func (r *TeamRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.Team, error) {
	e := r.db.getExec(ctx)

	const q = `
        SELECT id, name, created_at
        FROM teams
        WHERE id = $1
    `

	var t entity.Team
	err := e.QueryRowContext(ctx, q, id).Scan(
		&t.ID,
		&t.Name,
		&t.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return entity.Team{}, domain.ErrNotFound
	}
	if err != nil {
		return entity.Team{}, err
	}
	return t, nil
}
