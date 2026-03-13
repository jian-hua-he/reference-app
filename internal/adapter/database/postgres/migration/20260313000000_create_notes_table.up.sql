CREATE TABLE IF NOT EXISTS notes (
    id         TEXT PRIMARY KEY,
    text       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
