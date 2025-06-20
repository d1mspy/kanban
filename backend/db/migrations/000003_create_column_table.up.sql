CREATE TABLE IF NOT EXISTS "column"(
    id uuid PRIMARY KEY,
    board_id uuid REFERENCES "board"(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    name text NOT NULL,
    position smallint NOT NULL
);