package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/Desnn1ch/pr-reviewer-service/internal/app/service"
	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/common"
	dbinfra "github.com/Desnn1ch/pr-reviewer-service/internal/infrastructure/persistence/db"
	httpserver "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver"
	req "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/request"
	resp "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/response"
	"github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/handler"
)

var (
	baseURL    string
	httpSrv    *httptest.Server
	db         *dbinfra.DB
	pgC        testcontainers.Container
	cancelTest context.CancelFunc
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	cancelTest = cancel

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "app",
			"POSTGRES_PASSWORD": "app",
			"POSTGRES_DB":       "app",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(90 * time.Second),
	}
	var err error
	pgC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		panic(fmt.Errorf("start postgres container: %w", err))
	}

	host, err := pgC.Host(ctx)
	if err != nil {
		panic(fmt.Errorf("container host: %w", err))
	}
	mapped, err := pgC.MappedPort(ctx, "5432")
	if err != nil {
		panic(fmt.Errorf("mapped port: %w", err))
	}
	dsn := fmt.Sprintf("postgres://app:app@%s:%s/app?sslmode=disable", host, mapped.Port())

	wd, _ := os.Getwd()
	migrationsDir := filepath.Join(wd, "..", "migrations")
	cfg := dbinfra.Config{
		DSN:             dsn,
		MigrationsDir:   migrationsDir,
		MaxOpenConns:    5,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Minute,
	}
	db, err = dbinfra.New(ctx, cfg)
	if err != nil {
		_ = pgC.Terminate(context.Background())
		panic(fmt.Errorf("open db: %w", err))
	}

	repos := dbinfra.NewRepositories(db)
	teamSvc := service.NewTeamService(repos.Teams, repos.Users, repos.Tx)
	userSvc := service.NewUserService(repos.Users, repos.PRs)
	prSvc := service.NewPRService(repos.PRs, repos.Users, repos.Tx, common.StandardClock{})
	stSvc := service.NewStatsService(repos.PRs)

	teamH := handler.NewTeamHandler(teamSvc)
	userH := handler.NewUserHandler(userSvc)
	prH := handler.NewPRHandler(prSvc)
	statsH := handler.NewStatsHandler(stSvc)

	router := httpserver.NewRouter(teamH, userH, prH, statsH)
	httpSrv = httptest.NewServer(router)
	baseURL = httpSrv.URL

	code := m.Run()

	cancel()
	httpSrv.Close()
	_ = db.Close()
	_ = pgC.Terminate(context.Background())

	os.Exit(code)
}

func TestE2E_PRrev(t *testing.T) {
	teamName := "backend"
	authorID := uuid.New().String()
	reviewer1 := uuid.New().String()
	reviewer2 := uuid.New().String()
	reviewer3 := uuid.New().String()

	addReq := req.TeamAdd{
		TeamName: teamName,
		Members: []req.TeamMember{
			{UserID: authorID, Username: "author", IsActive: true},
			{UserID: reviewer1, Username: "u1", IsActive: true},
			{UserID: reviewer2, Username: "u2", IsActive: true},
			{UserID: reviewer3, Username: "u3", IsActive: true},
		},
	}

	var addResp resp.TeamAdd
	testPost(t, "/team/add", addReq, http.StatusOK, &addResp)
	if addResp.Team.TeamName != teamName {
		t.Fatalf("team name mismatch: got %q want %q", addResp.Team.TeamName, teamName)
	}
	if len(addResp.Team.Members) != 4 {
		t.Fatalf("members count: got %d want %d", len(addResp.Team.Members), 4)
	}

	prID := uuid.New().String()
	createReq := req.CreatePR{
		PullRequestID:   prID,
		PullRequestName: "Add search",
		AuthorID:        authorID,
	}
	var createResp resp.CreatePR
	testPost(t, "/pullRequest/create", createReq, http.StatusCreated, &createResp)

	if createResp.PR.PullRequestID != prID {
		t.Fatalf("pr id mismatch: got %q want %q", createResp.PR.PullRequestID, prID)
	}
	if createResp.PR.Status != "OPEN" {
		t.Fatalf("status after create: got %q want %q", createResp.PR.Status, "OPEN")
	}
	if got := len(createResp.PR.AssignedReviewers); got != 2 {
		t.Fatalf("assigned reviewers count: got %d want %d", got, 2)
	}
	for _, r := range createResp.PR.AssignedReviewers {
		if r == authorID {
			t.Fatalf("author must not be assigned as reviewer; got %q", r)
		}
	}

	reviewerToCheck := createResp.PR.AssignedReviewers[0]
	var reviews resp.UserReviews
	testGet(t, "/users/getReview?user_id="+reviewerToCheck, http.StatusOK, &reviews)
	if reviews.UserID != reviewerToCheck {
		t.Fatalf("reviews.user_id: got %q want %q", reviews.UserID, reviewerToCheck)
	}
	found := false
	for _, pr := range reviews.PullRequests {
		if pr.PullRequestID == prID && pr.Status == "OPEN" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("reviewer %s must have PR %s in OPEN status", reviewerToCheck, prID)
	}

	reassignReq := req.ReassignReviewer{
		PullRequestID: prID,
		OldUserID:     reviewerToCheck,
	}
	var reassignResp resp.ReassignReviewer
	testPost(t, "/pullRequest/reassign", reassignReq, http.StatusOK, &reassignResp)
	if reassignResp.ReplacedBy == reviewerToCheck {
		t.Fatalf("reassign must replace reviewer; got same id %q", reviewerToCheck)
	}
	if reassignResp.ReplacedBy == authorID {
		t.Fatalf("reassign must not assign author as reviewer")
	}
	if !contains(reassignResp.PR.AssignedReviewers, reassignResp.ReplacedBy) {
		t.Fatalf("new reviewer %s not in assigned list", reassignResp.ReplacedBy)
	}
	if contains(reassignResp.PR.AssignedReviewers, reviewerToCheck) {
		t.Fatalf("old reviewer %s still in assigned list", reviewerToCheck)
	}

	var mergeResp resp.MergePR
	testPost(t, "/pullRequest/merge", req.MergePR{PullRequestID: prID}, http.StatusOK, &mergeResp)
	if mergeResp.PR.Status != "MERGED" {
		t.Fatalf("status after merge: got %q want %q", mergeResp.PR.Status, "MERGED")
	}
	if mergeResp.PR.MergedAt == nil {
		t.Fatalf("mergedAt must be set after merge")
	}

	var errResp resp.Error
	testPost(t, "/pullRequest/reassign", reassignReq, http.StatusBadRequest, &errResp)
	if errResp.Error.Code != "PR_MERGED" {
		t.Fatalf("error code after reassign merged PR: got %q want %q", errResp.Error.Code, "PR_MERGED")
	}
}

func testPost(t *testing.T, path string, body any, wantStatus int, out any) {
	t.Helper()
	data, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, baseURL+path, bytes.NewReader(data))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode != wantStatus {
		t.Fatalf("status %s: got %d want %d", path, res.StatusCode, wantStatus)
	}
	if out != nil {
		if err := json.NewDecoder(res.Body).Decode(out); err != nil {
			t.Fatalf("decode response: %v", err)
		}
	}
}

func testGet(t *testing.T, path string, wantStatus int, out any) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, baseURL+path, nil)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("do request: %v", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode != wantStatus {
		t.Fatalf("status %s: got %d want %d", path, res.StatusCode, wantStatus)
	}
	if out != nil {
		if err := json.NewDecoder(res.Body).Decode(out); err != nil {
			t.Fatalf("decode response: %v", err)
		}
	}
}

func contains(ss []string, x string) bool {
	for _, s := range ss {
		if s == x {
			return true
		}
	}
	return false
}
