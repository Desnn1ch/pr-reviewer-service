package handler

import (
	"net/http"

	"github.com/Desnn1ch/pr-reviewer-service/internal/app/service"
	"github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/mapper"
)

type StatsHandler struct {
	svc *service.StatsService
}

func NewStatsHandler(svc *service.StatsService) *StatsHandler {
	return &StatsHandler{svc: svc}
}

func (h *StatsHandler) GetReviewerStats(w http.ResponseWriter, r *http.Request) {
	teamName := r.URL.Query().Get("team_name")
	if teamName == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "team_name is required")
		return
	}

	stats, err := h.svc.ReviewerStats(r.Context(), teamName)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		return
	}

	resp := mapper.ReviewerStatsToResponse(stats)
	writeJSON(w, http.StatusOK, resp)
}
