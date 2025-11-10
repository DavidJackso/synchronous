-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS session_leaderboard (
    session_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    tasks_completed INTEGER NOT NULL DEFAULT 0,
    focus_time INTEGER NOT NULL DEFAULT 0, -- в минутах
    score INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (session_id, user_id),
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_session_id ON session_leaderboard(session_id);
CREATE INDEX IF NOT EXISTS idx_user_id ON session_leaderboard(user_id);
CREATE INDEX IF NOT EXISTS idx_score ON session_leaderboard(score);

DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'update_updated_at_column') THEN
        CREATE TRIGGER update_session_leaderboard_updated_at BEFORE UPDATE ON session_leaderboard
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS session_leaderboard;
-- +goose StatementEnd

