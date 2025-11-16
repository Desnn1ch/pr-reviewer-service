package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/Desnn1ch/pr-reviewer-service/internal/app/service"
	"github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/mapper"
	req "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/request"
	resp "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/response"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) SetIsActive(w http.ResponseWriter, r *http.Request) {
	var body req.SetIsActive
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid json")
		return
	}
	if body.UserID == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "user_id is required")
		return
	}

	id, isActive, err := mapper.SetIsActiveRequestToArgs(body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid user_id")
		return
	}

	user, teamName, err := h.svc.SetActive(r.Context(), id, isActive)
	if err != nil {
		if handleDomainError(w, err) {
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		return
	}

	userResp := mapper.UserToResponse(user, teamName)
	writeJSON(w, http.StatusOK, resp.SetIsActive{User: userResp})
}

func (h *UserHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "user_id is required")
		return
	}

	id, err := uuid.Parse(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "BAD_REQUEST", "invalid user_id")
		return
	}

	prs, err := h.svc.GetReviews(r.Context(), id)
	if err != nil {
		if handleDomainError(w, err) {
			return
		}
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		return
	}

	respBody := resp.UserReviews{
		UserID:       id.String(),
		PullRequests: mapper.PRsToShortResponse(prs),
	}

	writeJSON(w, http.StatusOK, respBody)
}
