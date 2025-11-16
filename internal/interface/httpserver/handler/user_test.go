package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	respdto "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/response"
)

func TestUserHandler_SetIsActive_BadRequests(t *testing.T) {
	h := &UserHandler{svc: nil}

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "invalid JSON",
			body:       "{",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing user_id",
			body:       `{"is_active": true}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid user_id format",
			body:       `{"user_id": "not-a-uuid", "is_active": true}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			h.SetIsActive(w, req)

			res := w.Result()
			defer func() {
				_ = res.Body.Close()
			}()

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("status: got %d, want %d", res.StatusCode, tt.wantStatus)
			}

			var er respdto.Error
			if err := json.NewDecoder(res.Body).Decode(&er); err != nil {
				t.Fatalf("decode error response: %v", err)
			}
			if er.Error.Code != "BAD_REQUEST" {
				t.Errorf("error.code: got %q, want %q", er.Error.Code, "BAD_REQUEST")
			}
			if er.Error.Message == "" {
				t.Errorf("error.message should not be empty")
			}
		})
	}
}

func TestUserHandler_GetReview_BadRequests(t *testing.T) {
	h := &UserHandler{svc: nil}

	tests := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{
			name:       "missing user_id",
			url:        "/users/getReview",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid user_id format",
			url:        "/users/getReview?user_id=not-a-uuid",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()

			h.GetReview(w, req)

			res := w.Result()
			defer func() {
				_ = res.Body.Close()
			}()

			if res.StatusCode != tt.wantStatus {
				t.Fatalf("status: got %d, want %d", res.StatusCode, tt.wantStatus)
			}

			var er respdto.Error
			if err := json.NewDecoder(res.Body).Decode(&er); err != nil {
				t.Fatalf("decode error response: %v", err)
			}
			if er.Error.Code != "BAD_REQUEST" {
				t.Errorf("error.code: got %q, want %q", er.Error.Code, "BAD_REQUEST")
			}
		})
	}
}
