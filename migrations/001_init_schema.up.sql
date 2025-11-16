-- +goose Up
CREATE TABLE teams (
                       name TEXT PRIMARY KEY
);

CREATE TABLE users (
                        id         UUID PRIMARY KEY,
                        team_name  UUID NOT NULL REFERENCES teams(name) ON DELETE CASCADE,
                        name       TEXT NOT NULL,
                        is_active  BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE pull_requests (
                        id          UUID PRIMARY KEY,
                        title       TEXT NOT NULL,
                        author_id   UUID NOT NULL REFERENCES users(id),
                        status      TEXT NOT NULL,
                        created_at  TIMESTAMPTZ NOT NULL,
                        merged_at   TIMESTAMPTZ
);

CREATE TABLE pr_reviewers (
                        pr_id       UUID NOT NULL REFERENCES pull_requests(id) ON DELETE CASCADE,
                        reviewer_id UUID NOT NULL REFERENCES users(id),
                        PRIMARY KEY (pr_id, reviewer_id)
);
