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
	Tags          []int
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
	thread.Tags = getTags(db, thread.Id)

	return thread
}

// Return past 5 threads posted by user
func getRecentThreadsByUser(db *sql.DB, userid int) []Thread {
	sql_statement := `
	SELECT t.id
	FROM Threads t
	WHERE t.authorid = $1
	ORDER BY t.timestamp DESC
	LIMIT 5
	`
	rows, err := db.Query(sql_statement, userid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	recentThreads := []Thread{}
	for rows.Next() {
		var threadid string
		rows.Scan(&threadid)
		recentThreads = append(recentThreads, GetThreadById(db, threadid))
	}

	if rows.Err() != nil {
		panic(err)
	}

	return recentThreads
}

func getNumberOfThreadsByUser(db *sql.DB, userid int) int {
	sql_statement := `
	SELECT COUNT(*)
	FROM Threads t
	WHERE t.authorid = $1
	`
	rows, err := db.Query(sql_statement, userid)
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

func getNumberOfLikesToUserThreads(db *sql.DB, userid int) int {
	sql_statement := `
	SELECT t.id
	FROM Threads t
	WHERE t.authorid = $1
	`
	rows, err := db.Query(sql_statement, userid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var threadid int
		rows.Scan(&threadid)
		count += getLikesFromThreadId(db, threadid, true)
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

	comments := []int{}
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

func getTags(db *sql.DB, id int) []int {
	sql_statement := `
	SELECT tagId
	FROM Thread_Tags tt
	WHERE tt.threadId = $1
	`
	rows, err := db.Query(sql_statement, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	tags := []int{}
	for rows.Next() {
		var tag int
		err := rows.Scan(&tag)
		if err != nil {
			panic(err)
		}
		tags = append(tags, tag)
	}

	if rows.Err() != nil {
		panic(err)
	}

	return tags
}

func EditThreadById(db *sql.DB, title *string, content *string, tags []int, threadid int) error {

	tx, err := db.Begin()
	if err != nil {
		return errors.New("Unable to begin database transaction")
	}
	defer tx.Rollback()

	if title != nil {
		_, err := tx.Exec("UPDATE Threads SET title = $1 WHERE id = $2", title, threadid)
		if err != nil {
			return errors.New("Title has improper formatting")
		}
	}

	if content != nil {
		_, err := tx.Exec("UPDATE Threads SET content = $1 WHERE id = $2", content, threadid)
		if err != nil {
			return errors.New("Content has improper formatting")
		}
	}

	if tags != nil {
		_, err := tx.Exec("DELETE FROM Thread_Tags WHERE threadid = $1", threadid)
		if err != nil {
			return errors.New("Unable to delete thread tags")
		}

		for _, tagId := range tags {
			_, err := tx.Exec("INSERT INTO Thread_Tags VALUES ($1, $2)", threadid, tagId)
			if err != nil {
				return errors.New("Tags has improper formatting")
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("Unable to commit transaction")
	}

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
