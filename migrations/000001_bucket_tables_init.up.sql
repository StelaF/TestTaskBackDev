CREATE TABLE IF NOT EXISTS refresh_tokens
(
    user_id      UUID        NOT NULL,
    token_hash   TEXT        NOT NULL,
    expires_at   TIMESTAMPTZ  NOT NULL,
    PRIMARY KEY (user_ID)
);