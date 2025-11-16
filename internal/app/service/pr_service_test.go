package service

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/common"
	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/entity"
)

const teamName = "backend"

func TestPRService_Create_AssignsReviewers(t *testing.T) {
	ctx := context.Background()

	userRepo := newFakeUserRepo()
	prRepo := newFakePRRepo()
	tx := fakeTx{}
	clock := common.StandardClock{}

	authorID := uuid.New()
	r1 := uuid.New()
	r2 := uuid.New()
	r3 := uuid.New()

	userRepo.users[authorID] = entity.User{ID: authorID, TeamName: teamName, Name: "Author", IsActive: true}
	userRepo.users[r1] = entity.User{ID: r1, TeamName: teamName, Name: "R1", IsActive: true}
	userRepo.users[r2] = entity.User{ID: r2, TeamName: teamName, Name: "R2", IsActive: true}
	userRepo.users[r3] = entity.User{ID: r3, TeamName: teamName, Name: "R3", IsActive: true}

	svc := NewPRService(prRepo, userRepo, tx, clock)

	prID := uuid.New()
	pr, err := svc.Create(ctx, prID, "Add search", authorID)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if len(pr.Reviewers) != 2 {
		t.Fatalf("expected 2 reviewers, got %d", len(pr.Reviewers))
	}

	allowed := map[uuid.UUID]struct{}{
		r1: {},
		r2: {},
		r3: {},
	}

	for _, rid := range pr.Reviewers {
		if rid == authorID {
			t.Fatalf("author must not be a reviewer")
		}
		if _, ok := allowed[rid]; !ok {
			t.Fatalf("reviewer %s is not from active team members", rid)
		}
	}
}

func TestPRService_Create_NoOtherActiveUsers(t *testing.T) {
	ctx := context.Background()

	userRepo := newFakeUserRepo()
	prRepo := newFakePRRepo()
	tx := fakeTx{}
	clock := common.StandardClock{}

	authorID := uuid.New()

	userRepo.users[authorID] = entity.User{
		ID:       authorID,
		TeamName: teamName,
		Name:     "Author",
		IsActive: true,
	}

	svc := NewPRService(prRepo, userRepo, tx, clock)

	prID := uuid.New()
	pr, err := svc.Create(ctx, prID, "Lonely PR", authorID)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if len(pr.Reviewers) != 0 {
		t.Fatalf("expected 0 reviewers, got %d", len(pr.Reviewers))
	}
}

func TestPRService_Merge_Idempotent(t *testing.T) {
	ctx := context.Background()

	userRepo := newFakeUserRepo()
	prRepo := newFakePRRepo()
	tx := fakeTx{}
	clock := common.StandardClock{}

	svc := NewPRService(prRepo, userRepo, tx, clock)

	prID := uuid.New()
	prRepo.prs[prID] = entity.PR{
		ID:       prID,
		Title:    "Some PR",
		AuthorID: uuid.New(),
		Status:   entity.StatusOpen,
	}

	first, err := svc.Merge(ctx, prID)
	if err != nil {
		t.Fatalf("first Merge error: %v", err)
	}
	if !first.IsMerged() {
		t.Fatalf("first Merge: expected status MERGED")
	}
	if first.MergedAt == nil {
		t.Fatalf("first Merge: mergedAt is nil")
	}

	second, err := svc.Merge(ctx, prID)
	if err != nil {
		t.Fatalf("second Merge error: %v", err)
	}
	if !second.IsMerged() {
		t.Fatalf("second Merge: expected status MERGED")
	}
	if second.MergedAt == nil {
		t.Fatalf("second Merge: mergedAt is nil")
	}
	if !second.MergedAt.Equal(*first.MergedAt) {
		t.Fatalf("Merge must be idempotent: mergedAt changed (%v vs %v)", second.MergedAt, first.MergedAt)
	}
}

func TestPRService_Reassign_MergedPR(t *testing.T) {
	ctx := context.Background()

	userRepo := newFakeUserRepo()
	prRepo := newFakePRRepo()
	tx := fakeTx{}
	clock := common.StandardClock{}

	prID := uuid.New()
	oldID := uuid.New()

	prRepo.prs[prID] = entity.PR{
		ID:        prID,
		Title:     "PR",
		AuthorID:  uuid.New(),
		Status:    entity.StatusMerged,
		Reviewers: []uuid.UUID{oldID},
	}

	svc := NewPRService(prRepo, userRepo, tx, clock)

	_, err := svc.ReassignReviewer(ctx, prID, oldID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != common.ErrPRMerged {
		t.Fatalf("expected ErrPRMerged, got %v", err)
	}
}

func TestPRService_Reassign_OldReviewerNotAssigned(t *testing.T) {
	ctx := context.Background()

	userRepo := newFakeUserRepo()
	prRepo := newFakePRRepo()
	tx := fakeTx{}
	clock := common.StandardClock{}

	authorID := uuid.New()
	oldID := uuid.New()
	otherReviewer := uuid.New()

	userRepo.users[authorID] = entity.User{ID: authorID, TeamName: teamName, Name: "Author", IsActive: true}
	userRepo.users[oldID] = entity.User{ID: oldID, TeamName: teamName, Name: "Old", IsActive: true}
	userRepo.users[otherReviewer] = entity.User{ID: otherReviewer, TeamName: teamName, Name: "Other", IsActive: true}

	prID := uuid.New()
	prRepo.prs[prID] = entity.PR{
		ID:        prID,
		Title:     "PR",
		AuthorID:  authorID,
		Status:    entity.StatusOpen,
		Reviewers: []uuid.UUID{otherReviewer},
	}

	svc := NewPRService(prRepo, userRepo, tx, clock)

	_, err := svc.ReassignReviewer(ctx, prID, oldID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != common.ErrNotAssigned {
		t.Fatalf("expected ErrNotAssigned, got %v", err)
	}
}

func TestPRService_Reassign_NoCandidate(t *testing.T) {
	ctx := context.Background()

	userRepo := newFakeUserRepo()
	prRepo := newFakePRRepo()
	tx := fakeTx{}
	clock := common.StandardClock{}

	oldID := uuid.New()
	authorID := uuid.New()
	prID := uuid.New()

	userRepo.users[oldID] = entity.User{ID: oldID, TeamName: teamName, Name: "Old", IsActive: true}
	userRepo.users[authorID] = entity.User{ID: authorID, TeamName: teamName, Name: "Author", IsActive: true}

	prRepo.prs[prID] = entity.PR{
		ID:        prID,
		Title:     "PR",
		AuthorID:  authorID,
		Status:    entity.StatusOpen,
		Reviewers: []uuid.UUID{oldID},
	}

	svc := NewPRService(prRepo, userRepo, tx, clock)

	_, err := svc.ReassignReviewer(ctx, prID, oldID)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != common.ErrNoCandidate {
		t.Fatalf("expected ErrNoCandidate, got %v", err)
	}
}

func TestPRService_Reassign_Success(t *testing.T) {
	ctx := context.Background()

	userRepo := newFakeUserRepo()
	prRepo := newFakePRRepo()
	tx := fakeTx{}
	clock := common.StandardClock{}

	authorID := uuid.New()
	oldID := uuid.New()
	otherReviewer := uuid.New()
	newCandidate := uuid.New()

	userRepo.users[authorID] = entity.User{ID: authorID, TeamName: teamName, Name: "Author", IsActive: true}
	userRepo.users[oldID] = entity.User{ID: oldID, TeamName: teamName, Name: "Old", IsActive: true}
	userRepo.users[otherReviewer] = entity.User{ID: otherReviewer, TeamName: teamName, Name: "Other", IsActive: true}
	userRepo.users[newCandidate] = entity.User{ID: newCandidate, TeamName: teamName, Name: "New", IsActive: true}

	prID := uuid.New()
	prRepo.prs[prID] = entity.PR{
		ID:        prID,
		Title:     "PR",
		AuthorID:  authorID,
		Status:    entity.StatusOpen,
		Reviewers: []uuid.UUID{oldID, otherReviewer},
	}

	svc := NewPRService(prRepo, userRepo, tx, clock)

	res, err := svc.ReassignReviewer(ctx, prID, oldID)
	if err != nil {
		t.Fatalf("ReassignReviewer error: %v", err)
	}

	if len(res.Reviewers) != 2 {
		t.Fatalf("expected 2 reviewers, got %d", len(res.Reviewers))
	}

	foundOld := false
	foundNew := false
	for _, rid := range res.Reviewers {
		if rid == oldID {
			foundOld = true
		}
		if rid == newCandidate {
			foundNew = true
		}
	}

	if foundOld {
		t.Fatalf("old reviewer must be removed from reviewers")
	}
	if !foundNew {
		t.Fatalf("new candidate must be in reviewers")
	}
}
