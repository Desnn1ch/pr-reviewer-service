package db

import (
	"context"
	"database/sql"

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
		INSERT INTO teams (name)
		VALUES ($1)
	`

	_, err := e.ExecContext(ctx, q, team.Name)
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
		SELECT name
		FROM teams
		WHERE name = $1
	`

	var t entity.Team
	err := e.QueryRowContext(ctx, q, name).Scan(&t.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.Team{}, common.ErrNotFound
		}
		return entity.Team{}, err
	}

	return t, nil
}
