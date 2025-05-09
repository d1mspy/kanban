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

	QueryGetMaxColumnPosition = `SELECT COALESCE(MAX(position), 0) + 1 FROM "column" WHERE board_id = $1`

	QueryCheckBoardOwnership = `SELECT EXISTS (
		SELECT 1 FROM board WHERE id = $1 AND user_id = $2
	)`

	QueryGetColumnsCount = `SELECT COUNT(*) FROM "column" WHERE board_id = $1`

	QueryCreateColumn = `INSERT INTO "column" 
		(id, board_id, created_at, updated_at, name, position) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	QueryGetColumn = `SELECT "column".* FROM "column"
		JOIN board ON "column".board_id = board.id
		WHERE "column".id = $1 AND board.user_id = $2`

	QueryGetAllColumns = `SELECT "column".* FROM "column"
		JOIN board ON "column".board_id = board.id
		WHERE "column".board_id = $1
		AND board.user_id = $2`
	
	QueryDeleteColumn = `DELETE FROM "column" WHERE id = $1`
	
	QueryDecreaseColumnsPosition = `UPDATE "column" 
		SET position = position - 1
		WHERE board_id = $1
		AND position > $2`

	QueryGetBoardIDAndColumnPos =  `SELECT board_id, position FROM "column" WHERE id = $1`
	
	QueryMoveColumnsRight = `UPDATE "column"
		SET position = position + 1
		WHERE board_id = $1 AND position >= $2 AND position < $3`

	QueryMoveColumnsLeft = `UPDATE "column"
		SET position = position - 1
		WHERE board_id = $1 AND position > $2 AND position <= $3`

	QueryUpdateColumn = `UPDATE "column"
		SET name = COALESCE($1, name),
			position = COALESCE($2, position),
			updated_at = $3
		WHERE id = $4`
)