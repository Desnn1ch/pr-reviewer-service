package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Desnn1ch/pr-reviewer-service/internal/app/service"
	"github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/mapper"
	req "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/request"
	resp "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/response"
)

type PRHandler struct {
	svc *service.PRService
}

func NewPRHandler(svc *service.PRService) *PRHandler {
	return &PRHandler{svc: svc}
}

func (h *PRHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body req.CreatePR
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid json")
		return
	}
	if body.PullRequestID == "" || body.PullRequestName == "" || body.AuthorID == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "missing fields")
		return
	}

	id, title, authorID, err := mapper.CreatePRRequestToArgs(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid ids")
		return
	}

	pr, err := h.svc.Create(r.Context(), id, title, authorID)
	if err != nil {
		if handleDomainError(w, err) {
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		return
	}

	writeJSON(w, http.StatusCreated, resp.CreatePR{PR: mapper.PRToResponse(pr)})
}

func (h *PRHandler) Merge(w http.ResponseWriter, r *http.Request) {
	var body req.MergePR
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid json")
		return
	}
	if body.PullRequestID == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "pull_request_id is required")
		return
	}

	id, err := mapper.MergePRRequestToID(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid pull_request_id")
		return
	}

	pr, err := h.svc.Merge(r.Context(), id)
	if err != nil {
		if handleDomainError(w, err) {
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		return
	}

	writeJSON(w, http.StatusOK, resp.MergePR{PR: mapper.PRToResponse(pr)})
}

func (h *PRHandler) Reassign(w http.ResponseWriter, r *http.Request) {
	var body req.ReassignReviewer
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid json")
		return
	}
	if body.PullRequestID == "" || body.OldUserID == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "missing fields")
		return
	}

	prID, oldID, err := mapper.ReassignRequestToArgs(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid ids")
		return
	}

	pr, err := h.svc.ReassignReviewer(r.Context(), prID, oldID)
	if err != nil {
		if handleDomainError(w, err) {
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		return
	}

	replacedBy := ""
	for _, id := range pr.Reviewers {
		if id != oldID {
			replacedBy = id.String()
			break
		}
	}

	respBody := resp.ReassignReviewer{
		PR:         mapper.PRToResponse(pr),
		ReplacedBy: replacedBy,
	}

	writeJSON(w, http.StatusOK, respBody)
}
