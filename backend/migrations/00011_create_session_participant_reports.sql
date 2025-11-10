-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS session_participant_reports (
    report_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    tasks_completed INTEGER NOT NULL DEFAULT 0,
    focus_time INTEGER NOT NULL DEFAULT 0, -- в минутах
    PRIMARY KEY (report_id, user_id),
    FOREIGN KEY (report_id) REFERENCES session_reports(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_id ON session_participant_reports(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS session_participant_reports;
-- +goose StatementEnd

