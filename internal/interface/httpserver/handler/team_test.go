package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	respdto "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/response"
)

func TestTeamHandler_Add_BadRequests(t *testing.T) {
	h := &TeamHandler{svc: nil}

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
			name:       "missing team_name",
			body:       `{"members":[{"user_id":"c0f8a1c1-3a21-4b55-9e7c-4f8ba2e9d111","username":"Alice","is_active":true}]}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty members",
			body:       `{"team_name":"backend","members":[]}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid member user_id format",
			body:       `{"team_name":"backend","members":[{"user_id":"not-a-uuid","username":"Alice","is_active":true}]}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			h.Add(w, req)

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
			if er.Error.Code != BadRequestCode {
				t.Errorf("error.code: got %q, want %q", er.Error.Code, BadRequestCode)
			}
		})
	}
}

func TestTeamHandler_Get_BadRequests(t *testing.T) {
	h := &TeamHandler{svc: nil}

	tests := []struct {
		name       string
		url        string
		wantStatus int
	}{
		{
			name:       "missing team_name",
			url:        "/team/get",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty team_name",
			url:        "/team/get?team_name=",
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()

			h.Get(w, req)

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
			if er.Error.Code != BadRequestCode {
				t.Errorf("error.code: got %q, want %q", er.Error.Code, BadRequestCode)
			}
		})
	}
}
