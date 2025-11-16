package mapper

import (
	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/entity"
	resp "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/response"
)

func ReviewerStatsToResponse(stats []entity.ReviewerStats) resp.ReviewerStats {
	items := make([]resp.ReviewerStat, 0, len(stats))

	for _, s := range stats {
		items = append(items, resp.ReviewerStat{
			UserID:          s.UserID.String(),
			Username:        s.Username,
			TeamName:        s.TeamName,
			AssignedOpenPRs: s.AssignedOpenPRs,
		})
	}

	return resp.ReviewerStats{Items: items}
}
