CREATE TABLE IF NOT EXISTS "task"(
    id uuid PRIMARY KEY,
    column_id uuid REFERENCES "column"(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    position smallint NOT NULL,
    done boolean NOT NULL DEFAULT false,
    deadline timestamptz
);