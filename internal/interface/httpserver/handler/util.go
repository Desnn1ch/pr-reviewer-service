package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/common"
	resp "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/response"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, resp.Error{
		Error: resp.ErrorBody{
			Code:    code,
			Message: message,
		},
	})
}

func handleDomainError(w http.ResponseWriter, err error) bool {
	switch {
	case errors.Is(err, common.ErrTeamExists):
		writeError(w, http.StatusBadRequest, "TEAM_EXISTS", err.Error())
	case errors.Is(err, common.ErrPRExists):
		writeError(w, http.StatusConflict, "PR_EXISTS", err.Error())
	case errors.Is(err, common.ErrPRMerged):
		writeError(w, http.StatusBadRequest, "PR_MERGED", err.Error())
	case errors.Is(err, common.ErrNotAssigned):
		writeError(w, http.StatusBadRequest, "NOT_ASSIGNED", err.Error())
	case errors.Is(err, common.ErrNoCandidate):
		writeError(w, http.StatusBadRequest, "NO_CANDIDATE", err.Error())
	case errors.Is(err, common.ErrNotFound):
		writeError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
	case errors.Is(err, common.ErrUserInAnotherTeam):
		writeError(w, http.StatusBadRequest, "USER_IN_ANOTHER_TEAM", err.Error())
	default:
		return false
	}
	return true
}
