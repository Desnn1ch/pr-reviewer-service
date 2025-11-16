package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/common"
	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/entity"
)

func TestUserService_SetActive_ChangeState(t *testing.T) {
	ctx := context.Background()

	userRepo := newFakeUserRepo()
	prRepo := newFakePRRepo()

	id := uuid.New()
	userRepo.users[id] = entity.User{
		ID:       id,
		Name:     "Alice",
		IsActive: true,
	}

	svc := NewUserService(userRepo, prRepo)

	updated, err := svc.SetActive(ctx, id, false)
	if err != nil {
		t.Fatalf("SetActive returned error: %v", err)
	}

	if updated.IsActive != false {
		t.Fatalf("expected IsActive=false, got %v", updated.IsActive)
	}

	if userRepo.users[id].IsActive != false {
		t.Fatalf("repo user IsActive not updated")
	}

	if userRepo.setCalls != 1 {
		t.Fatalf("expected SetActive to be called once, got %d", userRepo.setCalls)
	}
}

func TestUserService_SetActive_NoChange(t *testing.T) {
	ctx := context.Background()

	userRepo := newFakeUserRepo()
	prRepo := newFakePRRepo()

	id := uuid.New()
	userRepo.users[id] = entity.User{
		ID:       id,
		Name:     "Alice",
		IsActive: true,
	}

	svc := NewUserService(userRepo, prRepo)

	updated, err := svc.SetActive(ctx, id, true)
	if err != nil {
		t.Fatalf("SetActive returned error: %v", err)
	}

	if updated.IsActive != true {
		t.Fatalf("expected IsActive=true, got %v", updated.IsActive)
	}

	if userRepo.setCalls != 0 {
		t.Fatalf("expected SetActive not to be called, got %d", userRepo.setCalls)
	}
}

func TestUserService_SetActive_GetByIDError(t *testing.T) {
	ctx := context.Background()

	userRepo := newFakeUserRepo()
	prRepo := newFakePRRepo()

	someErr := errors.New("db error")
	userRepo.getErr = someErr

	id := uuid.New()

	svc := NewUserService(userRepo, prRepo)

	_, err := svc.SetActive(ctx, id, false)
	if !errors.Is(err, someErr) {
		t.Fatalf("expected getErr (%v), got %v", someErr, err)
	}
}

func TestUserService_SetActive_SetActiveError(t *testing.T) {
	ctx := context.Background()

	userRepo := newFakeUserRepo()
	prRepo := newFakePRRepo()

	id := uuid.New()
	userRepo.users[id] = entity.User{
		ID:       id,
		Name:     "Alice",
		IsActive: true,
	}

	setErr := errors.New("update failed")
	userRepo.setErr = setErr

	svc := NewUserService(userRepo, prRepo)

	_, err := svc.SetActive(ctx, id, false)
	if !errors.Is(err, setErr) {
		t.Fatalf("expected setErr (%v), got %v", setErr, err)
	}
}

func TestUserService_GetReviews_Success(t *testing.T) {
	ctx := context.Background()

	userRepo := newFakeUserRepo()
	prRepo := newFakePRRepo()

	id := uuid.New()
	userRepo.users[id] = entity.User{
		ID:       id,
		Name:     "Alice",
		IsActive: true,
	}

	pr1 := entity.PR{ID: uuid.New(), Title: "PR1", Reviewers: []uuid.UUID{id}}
	pr2 := entity.PR{ID: uuid.New(), Title: "PR2", Reviewers: []uuid.UUID{id}}
	prRepo.prs[pr1.ID] = pr1
	prRepo.prs[pr2.ID] = pr2

	svc := NewUserService(userRepo, prRepo)

	prs, err := svc.GetReviews(ctx, id)
	if err != nil {
		t.Fatalf("GetReviews returned error: %v", err)
	}

	if len(prs) != 2 {
		t.Fatalf("expected 2 PRs, got %d", len(prs))
	}
}

func TestUserService_GetReviews_UserNotFound(t *testing.T) {
	ctx := context.Background()

	userRepo := newFakeUserRepo()
	prRepo := newFakePRRepo()

	svc := NewUserService(userRepo, prRepo)

	_, err := svc.GetReviews(ctx, uuid.New())
	if !errors.Is(err, common.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUserService_GetReviews_ListError(t *testing.T) {
	ctx := context.Background()

	userRepo := newFakeUserRepo()
	prRepo := newFakePRRepo()

	id := uuid.New()
	userRepo.users[id] = entity.User{
		ID:       id,
		Name:     "Alice",
		IsActive: true,
	}

	listErr := errors.New("list failed")
	prRepo.listErr = listErr

	svc := NewUserService(userRepo, prRepo)

	_, err := svc.GetReviews(ctx, id)
	if !errors.Is(err, listErr) {
		t.Fatalf("expected listErr (%v), got %v", listErr, err)
	}
}
