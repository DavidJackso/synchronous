-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS global_leaderboard (
    user_id VARCHAR(36) PRIMARY KEY,
    total_score INTEGER NOT NULL DEFAULT 0,
    total_sessions INTEGER NOT NULL DEFAULT 0,
    total_focus_time INTEGER NOT NULL DEFAULT 0, -- в минутах
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_total_score ON global_leaderboard(total_score);

DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'update_updated_at_column') THEN
        CREATE TRIGGER update_global_leaderboard_updated_at BEFORE UPDATE ON global_leaderboard
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS global_leaderboard;
-- +goose StatementEnd

