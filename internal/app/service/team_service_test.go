package service

import (
	"context"
	"errors"
	"github.com/Desnn1ch/pr-reviewer-service/internal/app"
	"testing"

	"github.com/google/uuid"

	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/common"
	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/entity"
)

func TestTeamService_CreateTeam_Success(t *testing.T) {
	ctx := context.Background()

	teamRepo := newFakeTeamRepo()
	userRepo := newFakeUserRepo()

	svc := NewTeamService(teamRepo, userRepo, fakeTx{})

	members := []entity.User{
		{ID: uuid.New(), Name: "Alice", IsActive: true},
		{ID: uuid.New(), Name: "Bob", IsActive: false},
	}
	teamName := "backend"

	team, users, err := svc.CreateTeam(ctx, teamName, members)
	if err != nil {
		t.Fatalf("CreateTeam returned error: %v", err)
	}

	if team.Name != teamName {
		t.Fatalf("expected team name %q, got %q", teamName, team.Name)
	}

	if len(users) != len(members) {
		t.Fatalf("expected %d users, got %d", len(members), len(users))
	}

	for i, u := range users {
		if u.ID != members[i].ID {
			t.Fatalf("user %d: expected ID %s, got %s", i, members[i].ID, u.ID)
		}
		if u.Name != members[i].Name {
			t.Fatalf("user %d: expected Name %q, got %q", i, members[i].Name, u.Name)
		}
		if u.IsActive != members[i].IsActive {
			t.Fatalf("user %d: expected IsActive %v, got %v", i, members[i].IsActive, u.IsActive)
		}
		if u.TeamName != teamName {
			t.Fatalf("user %d: expected TeamName %q, got %q", i, teamName, u.TeamName)
		}
	}

	if _, ok := teamRepo.teams[teamName]; !ok {
		t.Fatalf("team %q not stored in repo", teamName)
	}
}

func TestTeamService_CreateTeam_TeamExists(t *testing.T) {
	ctx := context.Background()

	teamRepo := newFakeTeamRepo()
	userRepo := newFakeUserRepo()

	existingName := "backend"
	teamRepo.teams[existingName] = entity.Team{Name: existingName}

	svc := NewTeamService(teamRepo, userRepo, fakeTx{})

	_, _, err := svc.CreateTeam(ctx, existingName, []entity.User{
		{ID: uuid.New(), Name: "Alice", IsActive: true},
	})
	if !errors.Is(err, common.ErrTeamExists) {
		t.Fatalf("expected ErrTeamExists, got %v", err)
	}
}

func TestTeamService_CreateTeam_UserInAnotherTeam(t *testing.T) {
	ctx := context.Background()

	teamRepo := newFakeTeamRepo()
	userRepo := newFakeUserRepo()

	existingID := uuid.New()
	userRepo.users[existingID] = entity.User{
		ID:       existingID,
		TeamName: "other-team",
		Name:     "Existing",
		IsActive: true,
	}

	svc := NewTeamService(teamRepo, userRepo, fakeTx{})

	members := []entity.User{
		{ID: existingID, Name: "Existing", IsActive: true},
	}

	_, _, err := svc.CreateTeam(ctx, "backend", members)
	if !errors.Is(err, common.ErrUserInAnotherTeam) {
		t.Fatalf("expected ErrUserInAnotherTeam, got %v", err)
	}

	if len(teamRepo.teams) != 0 {
		t.Fatalf("expected no teams to be created, got %d", len(teamRepo.teams))
	}
}

func TestTeamService_GetTeam_Success(t *testing.T) {
	ctx := context.Background()

	teamRepo := newFakeTeamRepo()
	userRepo := newFakeUserRepo()

	teamName := "backend"
	teamRepo.teams[teamName] = entity.Team{Name: teamName}

	u1 := entity.User{ID: uuid.New(), TeamName: teamName, Name: "Alice", IsActive: true}
	u2 := entity.User{ID: uuid.New(), TeamName: teamName, Name: "Bob", IsActive: false}
	uOther := entity.User{ID: uuid.New(), TeamName: "other", Name: "Charlie", IsActive: true}

	userRepo.users[u1.ID] = u1
	userRepo.users[u2.ID] = u2
	userRepo.users[uOther.ID] = uOther

	svc := NewTeamService(teamRepo, userRepo, fakeTx{})

	team, members, err := svc.GetTeam(ctx, teamName)
	if err != nil {
		t.Fatalf("GetTeam returned error: %v", err)
	}

	if team.Name != teamName {
		t.Fatalf("expected team name %q, got %q", teamName, team.Name)
	}

	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}

	for _, m := range members {
		if m.TeamName != teamName {
			t.Fatalf("member has wrong TeamName: %q", m.TeamName)
		}
	}
}

func TestTeamService_GetTeam_NotFound(t *testing.T) {
	ctx := context.Background()

	teamRepo := newFakeTeamRepo()
	userRepo := newFakeUserRepo()

	svc := NewTeamService(teamRepo, userRepo, fakeTx{})

	_, _, err := svc.GetTeam(ctx, "unknown")
	if !errors.Is(err, common.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
