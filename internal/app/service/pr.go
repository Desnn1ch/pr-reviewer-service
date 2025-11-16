package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Desnn1ch/pr-reviewer-service/internal/app"
	"github.com/Desnn1ch/pr-reviewer-service/internal/domain"
	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/entity"
)

type PRService struct {
	prs   app.PRRepo
	users app.UserRepo
	tx    app.TxManager
}

func NewPRService(prs app.PRRepo, users app.UserRepo, tx app.TxManager) *PRService {
	return &PRService{
		prs:   prs,
		users: users,
		tx:    tx,
	}
}

func (s *PRService) Create(ctx context.Context, id uuid.UUID, title string, authorID uuid.UUID) (entity.PR, error) {
	author, err := s.users.GetByID(ctx, authorID)
	if err != nil {
		return entity.PR{}, err
	}

	activeUsers, err := s.users.ListActiveByTeamID(ctx, author.TeamID)
	if err != nil {
		return entity.PR{}, err
	}

	candidates := make([]uuid.UUID, 0, len(activeUsers))
	for _, u := range activeUsers {
		if u.ID == author.ID {
			continue
		}
		candidates = append(candidates, u.ID)
	}

	reviewers := make([]uuid.UUID, 0, 2)
	for i := 0; i < len(candidates) && i < 2; i++ {
		reviewers = append(reviewers, candidates[i])
	}

	now := time.Now().UTC()

	pr := entity.PR{
		ID:        id,
		Title:     title,
		AuthorID:  authorID,
		Status:    entity.StatusOpen,
		CreatedAt: now,
		Reviewers: reviewers,
	}

	err = s.tx.InTx(ctx, func(txCtx context.Context) error {
		return s.prs.Create(txCtx, pr)
	})
	if err != nil {
		return entity.PR{}, err
	}

	return pr, nil
}

func (s *PRService) Merge(ctx context.Context, id uuid.UUID) (entity.PR, error) {
	var result entity.PR

	err := s.tx.InTx(ctx, func(txCtx context.Context) error {
		pr, err := s.prs.GetByID(txCtx, id)
		if err != nil {
			return err
		}

		if pr.IsMerged() {
			result = pr
			return nil
		}

		now := time.Now().UTC()
		pr.Status = entity.StatusMerged
		pr.MergedAt = &now

		if err := s.prs.Update(txCtx, pr); err != nil {
			return err
		}

		result = pr
		return nil
	})
	if err != nil {
		return entity.PR{}, err
	}

	return result, nil
}

func (s *PRService) ReassignReviewer(ctx context.Context, prID, oldReviewerID uuid.UUID) (entity.PR, error) {
	var result entity.PR

	err := s.tx.InTx(ctx, func(txCtx context.Context) error {
		pr, err := s.prs.GetByID(txCtx, prID)
		if err != nil {
			return err
		}

		if pr.IsMerged() {
			return domain.ErrPRMerged
		}

		idx := -1
		for i, r := range pr.Reviewers {
			if r == oldReviewerID {
				idx = i
				break
			}
		}
		if idx == -1 {
			return domain.ErrNotAssigned
		}

		oldReviewer, err := s.users.GetByID(txCtx, oldReviewerID)
		if err != nil {
			return err
		}

		activeUsers, err := s.users.ListActiveByTeamID(txCtx, oldReviewer.TeamID)
		if err != nil {
			return err
		}

		candidates := make([]uuid.UUID, 0, len(activeUsers))
		for _, u := range activeUsers {
			if u.ID == oldReviewer.ID {
				continue
			}
			if u.ID == pr.AuthorID {
				continue
			}
			isAlreadyReviewer := false
			for _, r := range pr.Reviewers {
				if r == u.ID {
					isAlreadyReviewer = true
					break
				}
			}
			if isAlreadyReviewer {
				continue
			}
			candidates = append(candidates, u.ID)
		}

		if len(candidates) == 0 {
			return domain.ErrNoCandidate
		}

		newReviewerID := candidates[0]
		pr.Reviewers[idx] = newReviewerID

		if err := s.prs.Update(txCtx, pr); err != nil {
			return err
		}

		result = pr
		return nil
	})
	if err != nil {
		return entity.PR{}, err
	}

	return result, nil
}
