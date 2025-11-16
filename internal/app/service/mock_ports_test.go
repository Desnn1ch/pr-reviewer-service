package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/common"
	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/entity"
)

type fakeTx struct{}

func (fakeTx) InTx(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

type fakeTeamRepo struct {
	teams map[string]entity.Team

	createErr error
	getErr    error
}

func newFakeTeamRepo() *fakeTeamRepo {
	return &fakeTeamRepo{
		teams: make(map[string]entity.Team),
	}
}

func (r *fakeTeamRepo) Create(ctx context.Context, team entity.Team) error {
	if r.createErr != nil {
		return r.createErr
	}

	if _, exists := r.teams[team.Name]; exists {
		return common.ErrTeamExists
	}

	r.teams[team.Name] = team
	return nil
}

func (r *fakeTeamRepo) GetByName(ctx context.Context, name string) (entity.Team, error) {
	if r.getErr != nil {
		return entity.Team{}, r.getErr
	}

	t, ok := r.teams[name]
	if !ok {
		return entity.Team{}, common.ErrNotFound
	}

	return t, nil
}

type fakeUserRepo struct {
	users map[uuid.UUID]entity.User

	getErr    error
	setErr    error
	upsertErr error

	setCalls int
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		users: make(map[uuid.UUID]entity.User),
	}
}

func (r *fakeUserRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	if r.getErr != nil {
		return entity.User{}, r.getErr
	}

	u, ok := r.users[id]
	if !ok {
		return entity.User{}, common.ErrNotFound
	}

	return u, nil
}

func (r *fakeUserRepo) ListByTeamName(ctx context.Context, teamName string) ([]entity.User, error) {
	var res []entity.User
	for _, u := range r.users {
		if u.TeamName == teamName {
			res = append(res, u)
		}
	}
	return res, nil
}

func (r *fakeUserRepo) ListActiveByTeamName(ctx context.Context, teamName string) ([]entity.User, error) {
	var res []entity.User
	for _, u := range r.users {
		if u.TeamName == teamName && u.IsActive {
			res = append(res, u)
		}
	}
	return res, nil
}

func (r *fakeUserRepo) UpsertMany(ctx context.Context, users []entity.User) error {
	if r.upsertErr != nil {
		return r.upsertErr
	}

	for _, u := range users {
		r.users[u.ID] = u
	}

	return nil
}

func (r *fakeUserRepo) SetActive(ctx context.Context, id uuid.UUID, active bool) error {
	r.setCalls++

	if r.setErr != nil {
		return r.setErr
	}

	u, ok := r.users[id]
	if !ok {
		return common.ErrNotFound
	}

	u.IsActive = active
	r.users[id] = u
	return nil
}

type fakePRRepo struct {
	prs map[uuid.UUID]entity.PR

	existsErr error
	createErr error
	getErr    error
	updateErr error
	listErr   error
}

func newFakePRRepo() *fakePRRepo {
	return &fakePRRepo{
		prs: make(map[uuid.UUID]entity.PR),
	}
}

func (r *fakePRRepo) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	if r.existsErr != nil {
		return false, r.existsErr
	}
	_, ok := r.prs[id]
	return ok, nil
}

func (r *fakePRRepo) Create(ctx context.Context, pr entity.PR) error {
	if r.createErr != nil {
		return r.createErr
	}

	if _, exists := r.prs[pr.ID]; exists {
		return common.ErrPRExists
	}

	r.prs[pr.ID] = pr
	return nil
}

func (r *fakePRRepo) GetByID(ctx context.Context, id uuid.UUID) (entity.PR, error) {
	if r.getErr != nil {
		return entity.PR{}, r.getErr
	}

	pr, ok := r.prs[id]
	if !ok {
		return entity.PR{}, common.ErrNotFound
	}

	return pr, nil
}

func (r *fakePRRepo) Update(ctx context.Context, pr entity.PR) error {
	if r.updateErr != nil {
		return r.updateErr
	}

	if _, ok := r.prs[pr.ID]; !ok {
		return common.ErrNotFound
	}

	r.prs[pr.ID] = pr
	return nil
}

func (r *fakePRRepo) ListByReviewerID(ctx context.Context, reviewerID uuid.UUID) ([]entity.PR, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}

	var res []entity.PR

	for _, pr := range r.prs {
		for _, rid := range pr.Reviewers {
			if rid == reviewerID {
				res = append(res, pr)
				break
			}
		}
	}

	return res, nil
}
