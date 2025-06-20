package postgres

const (

	// Auth queries

	QueryCreateUser = `
		INSERT INTO "user" 
		(id, username, email, hashed_password, created_at) 
		VALUES ($1, $2, $3, $4, $5)`

	QueryGetUserByEmail = `
		SELECT id, email, username, hashed_password 
		FROM "user" WHERE email=$1`

	// Board Queries

	QueryCreateBoard = `
		INSERT INTO board 
		(id, user_id, created_at, updated_at, name) 
		VALUES ($1, $2, $3, $4, $5)`

	QueryGetAllBoards = `
		SELECT * FROM board 
		WHERE user_id = $1 
		ORDER BY created_at`

	QueryGetBoard = `SELECT * FROM board WHERE id = $1`

	QueryUpdateBoard = `UPDATE board 
		SET updated_at = $1, name = $2
		WHERE id = $3`

	QueryDeleteBoard = `
		DELETE FROM board 
		WHERE id = $1`

	// Column queries

	QueryGetMaxColumnPosition = `
		SELECT COALESCE(MAX(position), 0) + 1 
		FROM "column" 
		WHERE board_id = $1`

	QueryGetColumnsCount = `
		SELECT COUNT(*) 
		FROM "column" 
		WHERE board_id = $1`

	QueryCreateColumn = `
		INSERT INTO "column" 
		(id, board_id, created_at, updated_at, name, position) 
		VALUES ($1, $2, $3, $4, $5, $6)`

	QueryGetColumn = `
		SELECT * 
		FROM "column"
		WHERE "column".id = $1`

	QueryGetAllColumns = `
		SELECT * 
		FROM "column"
		WHERE board_id = $1
		ORDER BY position`
	
	QueryDeleteColumn = `
		DELETE FROM "column" 
		WHERE id = $1`
	
	QueryDecreaseColumnsPosition = `
		UPDATE "column" 
		SET position = position - 1
		WHERE board_id = $1
		AND position > $2`

	QueryGetBoardIDAndColumnPos =  `
		SELECT board_id, position 
		FROM "column" 
		WHERE id = $1`
	
	QueryMoveColumnsRight = `
		UPDATE "column"
		SET position = position + 1
		WHERE board_id = $1 
		AND position >= $2 
		AND position < $3`

	QueryMoveColumnsLeft = `
		UPDATE "column"
		SET position = position - 1
		WHERE board_id = $1 
		AND position > $2 
		AND position <= $3`

	QueryUpdateColumn = `
		UPDATE "column"
		SET name = COALESCE($1, name),
			position = COALESCE($2, position),
			updated_at = $3
		WHERE id = $4`

	// Task queries

	QueryGetTasksCount = `
		SELECT COUNT(*) 
		FROM task 
		WHERE column_id = $1`

	QueryGetMaxTaskPosition = `
		SELECT COALESCE(MAX(position), 0) + 1 
		FROM task 
		WHERE column_id = $1`

	QueryCreateTask = `
		INSERT INTO task
		(id, column_id, created_at, updated_at, name, description, position, done)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	
	QueryGetAllTasks = `
		SELECT * FROM task 
		WHERE column_id = $1 
		ORDER BY position`

	QueryGetTask = `
		SELECT * 
		FROM task 
		WHERE id = $1`

	QueryUpdateTaskContent = `
		UPDATE task 
		SET name = COALESCE($1, name),
			description = COALESCE($2, description),
			done = COALESCE($3, done),
			deadline = COALESCE($4, deadline),
			updated_at = $5
		WHERE id = $6`
	
	QueryUpdateTaskColumn = `
		UPDATE task 
		SET column_id = $1,
			position = $2,
			updated_at = $3
		WHERE id = $4`

	QueryMoveTasksForInsert = `
		UPDATE task
		SET position = position + 1
		WHERE column_id = $1
		AND position >= $2`
	
	QueryMoveTaskForDelete = `
		UPDATE task
		SET position = position - 1
		WHERE column_id = $1
		AND position > $2`
	
	QueryGetColumnIDAndPosition = `
		SELECT column_id, position 
		FROM task 
		WHERE id = $1`

	QueryMoveTasksDown = `
		UPDATE task
		SET position = position + 1
		WHERE column_id = $1 
		AND position >= $2 
		AND position < $3`

	QueryMoveTasksUp = `
		UPDATE task
		SET position = position - 1
		WHERE column_id = $1 
		AND position > $2 
		AND position <= $3`

	QueryUpdateTaskPosition = `
		UPDATE task
		SET position = $1, 
			updated_at = $2
		WHERE id = $3`

	QueryDeleteTask = `
		DELETE FROM task 
		WHERE id = $1`

	// Queries for checking ownership

	QueryGetUserByBoardID = `
		SELECT user_id 
		FROM board 
		WHERE id = $1`

	QueryGetUserByColumnID = `
		SELECT board.user_id
		FROM board
		JOIN "column"
		ON "column".board_id = board.id
		WHERE "column".id = $1`
	
	QueryGetUserByTaskID = `
		SELECT board.user_id
		FROM task
		JOIN "column" ON task.column_id = "column".id
		JOIN board ON "column".board_id = board.id
		WHERE task.id = $1`

)