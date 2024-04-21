package sql

import (
	"database/sql"
)

func SetupProcedure(db *sql.DB) {
	procedureDir := "sql/procedure/"
	procedureFiles := []string{
		"make_verification.sql",
		"signup.sql",
		"update_all_comment_count.sql",
		"update_all_thread_count.sql",
		"verify_email.sql",
	}

	ExecSQLFiles(db, procedureDir, procedureFiles)
}
