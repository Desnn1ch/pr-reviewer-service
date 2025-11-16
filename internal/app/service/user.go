package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Desnn1ch/pr-reviewer-service/internal/app"
	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/entity"
)

type UserService struct {
	users app.UserRepo
	prs   app.PRRepo
}

func NewUserService(users app.UserRepo, prs app.PRRepo) *UserService {
	return &UserService{
		users: users,
		prs:   prs,
	}
}

func (s *UserService) SetActive(ctx context.Context, userID uuid.UUID, isActive bool) (entity.User, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return entity.User{}, err
	}

	if user.IsActive == isActive {
		return user, nil
	}

	if err := s.users.SetActive(ctx, userID, isActive); err != nil {
		return entity.User{}, err
	}

	user.IsActive = isActive
	return user, nil
}

func (s *UserService) GetReviews(ctx context.Context, userID uuid.UUID) ([]entity.PR, error) {
	_, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	prs, err := s.prs.ListByReviewerID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return prs, nil
}
