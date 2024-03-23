package sql

import (
	"database/sql"
)

func SetupTable(db *sql.DB) {
	tableDir := "sql/table/"
	tableFiles := []string{
		"users.sql",
		"follows.sql",
		"modules.sql",
		"threads.sql",
		"comments.sql",
		"subscribes.sql",
		"tokens.sql",
		"likes_comments.sql",
		"likes_threads.sql",
		"tags.sql",
		"thread_tags.sql",
		"reviews.sql",
		"bookmark_threads.sql",
		"email_verifications.sql",
	}

	ExecSQLFiles(db, tableDir, tableFiles)
}
