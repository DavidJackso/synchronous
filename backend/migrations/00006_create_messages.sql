-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS messages (
    id VARCHAR(36) PRIMARY KEY,
    session_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    text TEXT NOT NULL,
    max_message_id VARCHAR(255), -- ID сообщения в Max API
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_session_id ON messages(session_id);
CREATE INDEX IF NOT EXISTS idx_user_id ON messages(user_id);
CREATE INDEX IF NOT EXISTS idx_max_message_id ON messages(max_message_id);
CREATE INDEX IF NOT EXISTS idx_created_at ON messages(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
-- +goose StatementEnd

