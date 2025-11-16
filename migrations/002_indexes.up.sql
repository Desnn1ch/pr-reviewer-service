-- +goose Up
CREATE INDEX idx_pr_reviewers_reviewer_id ON pr_reviewers (reviewer_id);
CREATE INDEX idx_users_team_name_is_active ON users (team_name, is_active);
