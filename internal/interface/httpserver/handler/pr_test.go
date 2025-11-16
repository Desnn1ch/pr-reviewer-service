package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	respdto "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/response"
)

const BadRequestCode = "BAD_REQUEST"

func TestPRHandler_Create_BadRequests(t *testing.T) {
	h := &PRHandler{svc: nil}

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
			name:       "missing fields",
			body:       `{"pull_request_id":"pr-1","author_id":"11111111-1111-1111-1111-111111111111"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid pull_request_id uuid",
			body:       `{"pull_request_id":"not-a-uuid","pull_request_name":"Add search","author_id":"11111111-1111-1111-1111-111111111111"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid author_id uuid",
			body:       `{"pull_request_id":"11111111-1111-1111-1111-111111111111","pull_request_name":"Add search","author_id":"not-a-uuid"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			h.Create(w, req)

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

func TestPRHandler_Merge_BadRequests(t *testing.T) {
	h := &PRHandler{svc: nil}

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
			name:       "missing pull_request_id",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid pull_request_id uuid",
			body:       `{"pull_request_id":"not-a-uuid"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			h.Merge(w, req)

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

func TestPRHandler_Reassign_BadRequests(t *testing.T) {
	h := &PRHandler{svc: nil}

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
			name:       "missing fields",
			body:       `{"pull_request_id":"11111111-1111-1111-1111-111111111111"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid pull_request_id uuid",
			body:       `{"pull_request_id":"not-a-uuid","old_user_id":"11111111-1111-1111-1111-111111111111"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid old_user_id uuid",
			body:       `{"pull_request_id":"11111111-1111-1111-1111-111111111111","old_user_id":"not-a-uuid"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			h.Reassign(w, req)

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
