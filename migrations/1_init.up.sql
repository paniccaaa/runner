CREATE TABLE IF NOT EXISTS shared_codes (
    id SERIAL PRIMARY KEY,
    code TEXT NOT NULL,
    url TEXT NOT NULL,
    output TEXT,
    error TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

--"postgresql://postgres.xjbhwjugwklyfyloxvnj:bVYiPR8RyBVGLe62@aws-0-eu-central-1.pooler.supabase.com:6543/postgres"