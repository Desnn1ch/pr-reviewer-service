package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/common"
	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/entity"
)

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) UpsertMany(ctx context.Context, users []entity.User) error {
	if len(users) == 0 {
		return nil
	}

	e := r.db.getExec(ctx)

	const q = `
		INSERT INTO users (id, team_name, name, is_active)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET team_name = EXCLUDED.team_name,
			name = EXCLUDED.name,
			is_active = EXCLUDED.is_active;
	`

	for _, u := range users {
		_, err := e.ExecContext(ctx, q,
			u.ID,
			u.TeamName,
			u.Name,
			u.IsActive,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *UserRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	e := r.db.getExec(ctx)

	const q = `
		SELECT id, team_name, name, is_active
		FROM users
		WHERE id = $1
	`

	var u entity.User
	err := e.QueryRowContext(ctx, q, id).Scan(
		&u.ID,
		&u.TeamName,
		&u.Name,
		&u.IsActive,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, common.ErrNotFound
		}
		return entity.User{}, err
	}

	return u, nil
}

func (r *UserRepo) ListByTeamName(ctx context.Context, teamName string) ([]entity.User, error) {
	e := r.db.getExec(ctx)

	const q = `
		SELECT id, team_name, name, is_active
		FROM users
		WHERE team_name = $1
		ORDER BY name
	`

	rows, err := e.QueryContext(ctx, q, teamName)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var res []entity.User

	for rows.Next() {
		var u entity.User
		if err := rows.Scan(
			&u.ID,
			&u.TeamName,
			&u.Name,
			&u.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *UserRepo) ListActiveByTeamName(ctx context.Context, teamName string) ([]entity.User, error) {
	e := r.db.getExec(ctx)

	const q = `
		SELECT id, team_name, name, is_active
		FROM users
		WHERE team_name = $1
		  AND is_active = TRUE
		ORDER BY name
	`

	rows, err := e.QueryContext(ctx, q, teamName)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var res []entity.User

	for rows.Next() {
		var u entity.User
		if err := rows.Scan(
			&u.ID,
			&u.TeamName,
			&u.Name,
			&u.IsActive,
		); err != nil {
			return nil, err
		}
		res = append(res, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *UserRepo) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	e := r.db.getExec(ctx)

	const q = `
		UPDATE users
		SET is_active = $2
		WHERE id = $1
	`

	res, err := e.ExecContext(ctx, q, id, active)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return common.ErrNotFound
	}

	return nil
}
