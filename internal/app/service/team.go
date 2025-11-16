package service

import (
	"context"
	"errors"
	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/common"

	"github.com/google/uuid"

	"github.com/Desnn1ch/pr-reviewer-service/internal/app"
	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/entity"
)

type TeamService struct {
	teams app.TeamRepo
	users app.UserRepo
	tx    app.TxManager
	clock common.Clock
}

func NewTeamService(teams app.TeamRepo, users app.UserRepo, tx app.TxManager, clock common.Clock) *TeamService {
	return &TeamService{
		teams: teams,
		users: users,
		tx:    tx,
		clock: clock,
	}
}

func (s *TeamService) CreateTeam(ctx context.Context, name string, members []entity.User) (entity.Team, []entity.User, error) {
	_, err := s.teams.GetByName(ctx, name)
	if err == nil {
		return entity.Team{}, nil, common.ErrTeamExists
	}
	if err != nil && err != common.ErrNotFound {
		return entity.Team{}, nil, err
	}

	now := s.clock.Now()

	team := entity.Team{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: now,
	}

	users := make([]entity.User, len(members))
	for i, m := range members {
		u := m
		u.TeamID = team.ID
		if u.CreatedAt.IsZero() {
			u.CreatedAt = now
		}
		users[i] = u
	}

	err = s.tx.InTx(ctx, func(txCtx context.Context) error {
		for _, u := range users {
			existing, err := s.users.GetByID(txCtx, u.ID)
			if err != nil && !errors.Is(err, common.ErrNotFound) {
				return err
			}
			if err == nil && existing.TeamID != team.ID {
				return common.ErrUserInAnotherTeam
			}
		}

		if err := s.teams.Create(txCtx, team); err != nil {
			return err
		}
		if len(users) > 0 {
			if err := s.users.UpsertMany(txCtx, users); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return entity.Team{}, nil, err
	}

	return team, users, nil
}

func (s *TeamService) GetTeam(ctx context.Context, name string) (entity.Team, []entity.User, error) {
	team, err := s.teams.GetByName(ctx, name)
	if err != nil {
		return entity.Team{}, nil, err
	}

	members, err := s.users.ListByTeamID(ctx, team.ID)
	if err != nil {
		return entity.Team{}, nil, err
	}

	return team, members, nil
}
