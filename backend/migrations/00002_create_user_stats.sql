-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_stats (
    user_id VARCHAR(36) PRIMARY KEY,
    total_sessions INTEGER NOT NULL DEFAULT 0,
    total_focus_time INTEGER NOT NULL DEFAULT 0, -- в минутах
    current_streak INTEGER NOT NULL DEFAULT 0, -- в днях
    last_session_date DATE,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'update_updated_at_column') THEN
        CREATE FUNCTION update_updated_at_column()
        RETURNS TRIGGER AS $$
        BEGIN
            NEW.updated_at = CURRENT_TIMESTAMP;
            RETURN NEW;
        END;
        $$ language 'plpgsql';
    END IF;
END $$;

DO $$ 
BEGIN
    IF EXISTS (SELECT 1 FROM pg_proc WHERE proname = 'update_updated_at_column') THEN
        CREATE TRIGGER update_user_stats_updated_at BEFORE UPDATE ON user_stats
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_stats;
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;
-- +goose StatementEnd

