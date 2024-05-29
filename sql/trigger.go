package sql

import (
	"database/sql"
)

func SetupTrigger(db *sql.DB) {
	triggerDir := "sql/trigger/"
	triggerFiles := []string{
		"comment_comments_count_update_trigger.sql",
		"comment_likes_count_update_trigger.sql",
		"thread_comments_count_update_trigger.sql",
		"thread_likes_count_update_trigger.sql",
	}

	ExecSQLFiles(db, triggerDir, triggerFiles)
}
