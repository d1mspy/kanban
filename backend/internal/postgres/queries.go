package postgres

const (
	QueryCreateUserTable = `CREATE TABLE IF NOT EXISTS "user"(
		id uuid PRIMARY KEY,
		created_at timestamptz NOT NULL,
		username text NOT NULL UNIQUE,
		hashed_password text NOT NULL
	)`

	QueryCreateBoardTable = `CREATE TABLE IF NOT EXISTS board(
		id uuid PRIMARY KEY,
		user_id uuid REFERENCES "user"(id) ON DELETE CASCADE,
		created_at timestamptz NOT NULL,
		updated_at timestamptz NOT NULL,
		name text NOT NULL
	)`

	QueryCreateColumnTable = `CREATE TABLE IF NOT EXISTS "column"(
		id uuid PRIMARY KEY,
		board_id uuid REFERENCES "board"(id),
		created_at timestamptz NOT NULL,
		updated_at timestamptz NOT NULL,
		name text NOT NULL,
		position smallint NOT NULL
	)`

	QueryCreateTaskTable = `CREATE TABLE IF NOT EXISTS task(
		id uuid PRIMARY KEY,
		column_id uuid REFERENCES "column"(id),
		created_at timestamptz NOT NULL,
		updated_at timestamptz NOT NULL,
		name text NOT NULL,
		description text NOT NULL,
		position smallint NOT NULL
	)`

	QueryCreateUser = `INSERT INTO "user" (id, username, hashed_password, created_at) VALUES ($1, $2, $3, $4)`

	QueryGetUserByUsername = `SELECT id, username, hashed_password FROM "user" WHERE username=$1`

	QueryCreateBoard = `INSERT INTO board (id, user_id, created_at, updated_at, name) VALUES ($1, $2, $3, $4, $5)`

	QueryGetAllBoards = `SELECT * FROM board WHERE user_id = $1 ORDER BY created_at`

	QueryGetBoard = `SELECT * FROM board WHERE id = $1 AND user_id = $2`

	QueryUpdateBoard = `UPDATE board 
		SET updated_at = $1, name = $2
		WHERE id = $3 AND user_id = $4`

	QueryDeleteBoard = `DELETE FROM board WHERE id = $1 AND user_id = $2`
)