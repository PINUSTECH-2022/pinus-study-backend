package sql

import (
	"database/sql"
)

func SetupProcedure(db *sql.DB) {
	procedureDir := "sql/procedure/"
	procedureFiles := []string{
		"signup.sql",
		"make_verification.sql",
		"verify_email.sql",
		"update_all_thread_count.sql",
	}

	ExecSQLFiles(db, procedureDir, procedureFiles)
}
