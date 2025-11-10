-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS session_reports (
    id VARCHAR(36) PRIMARY KEY,
    session_id VARCHAR(36) NOT NULL UNIQUE,
    tasks_completed INTEGER NOT NULL DEFAULT 0,
    tasks_total INTEGER NOT NULL DEFAULT 0,
    focus_time INTEGER NOT NULL DEFAULT 0, -- в минутах
    break_time INTEGER NOT NULL DEFAULT 0, -- в минутах
    cycles_completed INTEGER NOT NULL DEFAULT 0,
    completed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_session_id ON session_reports(session_id);
CREATE INDEX IF NOT EXISTS idx_completed_at ON session_reports(completed_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS session_reports;
-- +goose StatementEnd

