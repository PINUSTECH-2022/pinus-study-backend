package database

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

type Thread struct {
	Id            int
	Title         string
	Content       string
	AuthorId      int
	Username      string
	Timestamp     string
	ModuleId      string
	LikesCount    int
	DislikesCount int
	Comments      []int
}

func GetThreadById(db *sql.DB, threadid string) Thread {
	rows, err := db.Query("SELECT * FROM Threads WHERE id = $1", threadid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var thread Thread

	for rows.Next() {
		err := rows.Scan(&thread.Id, &thread.Title, &thread.Content, &thread.ModuleId, &thread.AuthorId, &thread.Timestamp)
		if err != nil {
			panic(err)
		}
	}

	if rows.Err() != nil {
		panic(err)
	}

	rows, err = db.Query("SELECT username FROM Users WHERE id = $1", thread.AuthorId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&thread.Username)
		if err != nil {
			panic(err)
		}
	}

	if rows.Err() != nil {
		panic(err)
	}

	thread.LikesCount = getLikesFromThreadId(db, thread.Id, true)
	thread.DislikesCount = getLikesFromThreadId(db, thread.Id, false)
	thread.Comments = getComments(db, thread.Id)

	return thread
}

func getLikesFromThreadId(db *sql.DB, id int, status bool) int {
	sql_statement := `
	SELECT COUNT(*)
	FROM Likes_Threads
	WHERE state = $1 AND threadid = $2
	`
	rows, err := db.Query(sql_statement, status, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			panic(err)
		}
	}

	if rows.Err() != nil {
		panic(err)
	}

	return count
}

func getComments(db *sql.DB, id int) []int {
	sql_statement := `
	SELECT c.id
	FROM Comments c
	WHERE c.threadid = $1
	`
	rows, err := db.Query(sql_statement, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var comments []int
	for rows.Next() {
		var comment int
		err := rows.Scan(&comment)
		if err != nil {
			panic(err)
		}
		comments = append(comments, comment)
	}

	if rows.Err() != nil {
		panic(err)
	}

	return comments
}

func EditThreadById(db *sql.DB, title string, content string, threadid int) error {
	rows, err := db.Query("UPDATE Threads SET title = $1, content = $2 WHERE id = $3 ", title, content, threadid)
	if err != nil {
		return errors.New(err.Error())
	}
	defer rows.Close()
	return nil
}

func PostComment(db *sql.DB, authorid int, content string, parentid int, threadid int) error {
	rows, err := db.Query("INSERT INTO Comments (authorid, content, id, is_deleted, parentid, threadid, timestamp) VALUES ($1, $2, $3, $4, $5, $6, $7)", authorid, content, getCommentId(db), false, parentid, threadid, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return errors.New(err.Error())
	}
	defer rows.Close()
	return nil
}

func getCommentId(db *sql.DB) int {
	sql_statement := `
	SELECT COUNT(*)
	FROM Comments
	`
	rows, err := db.Query(sql_statement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			panic(err)
		}
	}

	if rows.Err() != nil {
		panic(err)
	}

	return count + 1
}
