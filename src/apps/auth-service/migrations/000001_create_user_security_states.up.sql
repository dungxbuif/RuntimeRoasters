CREATE TABLE user_security_states (
    subject             TEXT PRIMARY KEY,
    tokens_valid_after  TIMESTAMPTZ NOT NULL DEFAULT '1970-01-01 00:00:00+00',
    reason              TEXT NOT NULL DEFAULT '',
    updated_by          TEXT NOT NULL,
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version             BIGINT NOT NULL DEFAULT 1
);
