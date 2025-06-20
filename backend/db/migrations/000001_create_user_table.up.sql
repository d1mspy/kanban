CREATE TABLE IF NOT EXISTS "user"(
	id uuid PRIMARY KEY,
	created_at timestamptz NOT NULL,
    email TEXT NOT NULL UNIQUE,
	username text NOT NULL,
	hashed_password text NOT NULL
);