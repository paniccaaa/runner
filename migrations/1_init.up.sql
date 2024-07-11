CREATE TABLE IF NOT EXISTS shared_codes (
    id SERIAL PRIMARY KEY,
    code TEXT NOT NULL,
    output TEXT,
    error TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
