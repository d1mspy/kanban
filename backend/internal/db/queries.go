package db

const (
	QueryCreateUserTable = `CREATE TABLE IF NOT EXISTS "user"(
		id uuid PRIMARY KEY,
		created_at timestamptz NOT NULL,
		username text NOT NULL UNIQUE,
		hashed_password text NOT NULL
	)`
)