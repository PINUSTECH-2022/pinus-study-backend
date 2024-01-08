package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func BookmarkThread(db *sql.DB, threadId int, userId int) error {
	sql_statement := `
	INSERT INTO bookmark_threads(thread_id, user_id)
	VALUES ($1, $2);
	`
	_, err := db.Exec(sql_statement, threadId, userId)
	return err
}
