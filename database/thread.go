package database

import (
	"database/sql"

	"time"

	_ "github.com/lib/pq"
)

type Thread struct {
	Id        string
	Title     string
	Content   string
	AuthorId  string
	Timestamp time.Time
	ModuleId  string
}

func GetThreads(db *sql.DB, id string) []Thread {
	rows, err := db.Query("SELECT * FROM Threads WHERE id = $1", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var threads []Thread

	for rows.Next() {
		var thread Thread
		err := rows.Scan(&thread.Id, &thread.Title, &thread.Content, &thread.AuthorId, &thread.Timestamp, &thread.ModuleId)
		if err != nil {
			panic(err)
		}
		threads = append(threads, thread)
	}

	if rows.Err() != nil {
		panic(err)
	}

	return threads
}
