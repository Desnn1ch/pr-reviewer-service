package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Desnn1ch/pr-reviewer-service/internal/app/service"
	"github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/mapper"
	req "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/request"
	resp "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/response"
)

type TeamHandler struct {
	svc *service.TeamService
}

func NewTeamHandler(svc *service.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) Add(w http.ResponseWriter, r *http.Request) {
	var body req.TeamAdd
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid json")
		return
	}
	if body.TeamName == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "team_name is required")
		return
	}

	name, members, err := mapper.TeamAddRequestToArgs(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid team members")
		return
	}

	team, users, err := h.svc.CreateTeam(r.Context(), name, members)
	if err != nil {
		if handleDomainError(w, err) {
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		return
	}

	teamResp := mapper.TeamToResponse(team, users)
	writeJSON(w, http.StatusOK, resp.TeamAdd{Team: teamResp})
}

func (h *TeamHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("team_name")
	if name == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "team_name is required")
		return
	}

	team, users, err := h.svc.GetTeam(r.Context(), name)
	if err != nil {
		if handleDomainError(w, err) {
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		return
	}

	teamResp := mapper.TeamToResponse(team, users)
	writeJSON(w, http.StatusOK, teamResp)
}
