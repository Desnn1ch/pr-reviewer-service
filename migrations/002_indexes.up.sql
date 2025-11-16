-- +goose Up
CREATE UNIQUE INDEX ux_teams_name ON teams (name);
CREATE INDEX idx_users_team_id ON users (team_id);
CREATE INDEX idx_pr_reviewers_reviewer_id ON pr_reviewers (reviewer_id);
CREATE INDEX idx_users_team_id_is_active ON users (team_id, is_active);
