package postgres

const (
	QueryCreateUserTable = `CREATE TABLE IF NOT EXISTS "user"(
		id uuid PRIMARY KEY,
		created_at timestamptz NOT NULL,
		username text NOT NULL UNIQUE,
		hashed_password text NOT NULL
	)`

	QueryCreateUser = `INSERT INTO "user" (id, username, hashed_password, created_at) VALUES ($1, $2, $3, $4)`

	QueryGetUserByUsername = `SELECT id, username, hashed_password FROM "user" WHERE username=$1`
)