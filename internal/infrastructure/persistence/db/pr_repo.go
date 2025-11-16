package db

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/common"
	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/entity"
)

type PRRepo struct {
	db *DB
}

func NewPRRepo(db *DB) *PRRepo {
	return &PRRepo{db: db}
}

func (r *PRRepo) Create(ctx context.Context, pr entity.PR) error {
	e := r.db.getExec(ctx)

	const qPR = `
		INSERT INTO pull_requests (id, title, author_id, status, created_at, merged_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := e.ExecContext(ctx, qPR,
		pr.ID,
		pr.Title,
		pr.AuthorID,
		string(pr.Status),
		pr.CreatedAt,
		pr.MergedAt,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return common.ErrPRExists
		}
		return err
	}

	if len(pr.Reviewers) == 0 {
		return nil
	}

	const insertReviewer = `
		INSERT INTO pr_reviewers (pr_id, reviewer_id)
		VALUES ($1, $2)
	`

	for _, reviewerID := range pr.Reviewers {
		if _, err := e.ExecContext(ctx, insertReviewer, pr.ID, reviewerID); err != nil {
			return err
		}
	}

	return nil
}

func (r *PRRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.PR, error) {
	q := r.db.getExec(ctx)

	const qPR = `
		SELECT id, title, author_id, status, created_at, merged_at
		FROM pull_requests
		WHERE id = $1
	`

	var pr entity.PR
	var status string

	err := q.QueryRowContext(ctx, qPR, id).Scan(
		&pr.ID,
		&pr.Title,
		&pr.AuthorID,
		&status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.PR{}, common.ErrNotFound
		}
		return entity.PR{}, err
	}

	pr.Status = entity.PRStatus(status)

	reviewers, err := r.loadReviewers(ctx, id)
	if err != nil {
		return entity.PR{}, err
	}
	pr.Reviewers = reviewers

	return pr, nil
}

func (r *PRRepo) Update(ctx context.Context, pr entity.PR) error {
	e := r.db.getExec(ctx)

	const qPR = `
		UPDATE pull_requests
		SET title = $2,
			author_id = $3,
			status = $4,
			merged_at = $5
		WHERE id = $1
	`

	res, err := e.ExecContext(ctx, qPR,
		pr.ID,
		pr.Title,
		pr.AuthorID,
		string(pr.Status),
		pr.MergedAt,
	)
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

	const qDel = `DELETE FROM pr_reviewers WHERE pr_id = $1`

	if _, err := e.ExecContext(ctx, qDel, pr.ID); err != nil {
		return err
	}

	if len(pr.Reviewers) == 0 {
		return nil
	}

	const qIns = `
		INSERT INTO pr_reviewers (pr_id, reviewer_id)
		VALUES ($1, $2)
	`

	for _, reviewerID := range pr.Reviewers {
		if _, err := e.ExecContext(ctx, qIns, pr.ID, reviewerID); err != nil {
			return err
		}
	}

	return nil
}

func (r *PRRepo) ListByReviewerID(ctx context.Context, reviewerID uuid.UUID) ([]entity.PR, error) {
	q := r.db.getExec(ctx)

	const query = `
		SELECT pr.id, pr.title, pr.author_id, pr.status, pr.created_at, pr.merged_at
		FROM pull_requests pr
		JOIN pr_reviewers r ON r.pr_id = pr.id
		WHERE r.reviewer_id = $1
		ORDER BY pr.created_at DESC
	`

	rows, err := q.QueryContext(ctx, query, reviewerID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("rows close error: %v", cerr)
		}
	}()

	var result []entity.PR

	for rows.Next() {
		var pr entity.PR
		var status string

		if err := rows.Scan(
			&pr.ID,
			&pr.Title,
			&pr.AuthorID,
			&status,
			&pr.CreatedAt,
			&pr.MergedAt,
		); err != nil {
			return nil, err
		}

		pr.Status = entity.PRStatus(status)

		reviewers, err := r.loadReviewers(ctx, pr.ID)
		if err != nil {
			return nil, err
		}
		pr.Reviewers = reviewers

		result = append(result, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *PRRepo) loadReviewers(ctx context.Context, prID uuid.UUID) ([]uuid.UUID, error) {
	q := r.db.getExec(ctx)

	const query = `
		SELECT reviewer_id
		FROM pr_reviewers
		WHERE pr_id = $1
	`

	rows, err := q.QueryContext(ctx, query, prID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			log.Printf("rows close error: %v", cerr)
		}
	}()

	var reviewers []uuid.UUID

	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		reviewers = append(reviewers, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reviewers, nil
}

func (r *PRRepo) ListReviewerStats(ctx context.Context, teamName string) ([]entity.ReviewerStats, error) {
	e := r.db.getExec(ctx)

	const q = `
		SELECT
			u.id,
			u.name,
			u.team_name,
			COUNT(*) AS assigned_open_prs
		FROM pr_reviewers prr
		JOIN pull_requests p ON p.id = prr.pr_id
		JOIN users u         ON u.id = prr.reviewer_id
		WHERE p.status = 'OPEN'
		  AND u.team_name = $1
		GROUP BY u.id, u.name, u.team_name
		ORDER BY assigned_open_prs DESC
	`

	rows, err := e.QueryContext(ctx, q, teamName)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var res []entity.ReviewerStats

	for rows.Next() {
		var s entity.ReviewerStats
		if err := rows.Scan(
			&s.UserID,
			&s.Username,
			&s.TeamName,
			&s.AssignedOpenPRs,
		); err != nil {
			return nil, err
		}
		res = append(res, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}
