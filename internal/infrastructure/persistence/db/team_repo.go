package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/common"
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
	INSERT INTO teams (id, name)
	VALUES ($1, $2)
	`

	_, err := e.ExecContext(ctx, q,
		team.ID,
		team.Name,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return common.ErrTeamExists
		}
		return err
	}

	return nil
}

func (r *TeamRepo) GetByName(ctx context.Context, name string) (entity.Team, error) {
	e := r.db.getExec(ctx)

	const q = `
	SELECT id, name
	FROM teams
	WHERE name = $1
	`

	var t entity.Team
	err := e.QueryRowContext(ctx, q, name).Scan(
		&t.ID,
		&t.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Team{}, common.ErrNotFound
		}
		return entity.Team{}, err
	}

	return t, nil
}

func (r *TeamRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.Team, error) {
	e := r.db.getExec(ctx)

	const q = `
	SELECT id, name
	FROM teams
	WHERE id = $1
	`

	var t entity.Team
	err := e.QueryRowContext(ctx, q, id).Scan(
		&t.ID,
		&t.Name,
	)
	if err == sql.ErrNoRows {
		return entity.Team{}, common.ErrNotFound
	}
	if err != nil {
		return entity.Team{}, err
	}

	return t, nil
}
