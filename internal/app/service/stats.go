package service

import (
	"context"

	"github.com/Desnn1ch/pr-reviewer-service/internal/app"
	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/entity"
)

type StatsService struct {
	prs app.PRRepo
}

func NewStatsService(prs app.PRRepo) *StatsService {
	return &StatsService{prs: prs}
}

func (s *StatsService) ReviewerStats(ctx context.Context, teamName string) ([]entity.ReviewerStats, error) {
	return s.prs.ListReviewerStats(ctx, teamName)
}
