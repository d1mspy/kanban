CREATE TABLE IF NOT EXISTS "board"(
    id uuid PRIMARY KEY,
    user_id uuid REFERENCES "user"(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    name text NOT NULL
);