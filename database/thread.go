package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
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
	IsDeleted     bool
	Comments      []int
	Tags          []int
}

func GetThreadById(db *sql.DB, threadid string) Thread {
	fmt.Println("GetThreadById")
	fmt.Println(threadid)
	threadidInt, err := strconv.Atoi(threadid)
	rows, err := db.Query("SELECT id, title, content, moduleid, authorid, timestamp FROM Threads WHERE id = $1 AND is_deleted = 'f'", threadidInt)

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer rows.Close()
	fmt.Println("Query done")

	var thread Thread
	fmt.Println("Initial: ", thread)
	for rows.Next() {
		err := rows.Scan(&thread.Id, &thread.Title, &thread.Content, &thread.ModuleId, &thread.AuthorId, &thread.Timestamp, &thread.IsDeleted)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(thread)
	if rows.Err() != nil {
		panic(err)
	}
	fmt.Println("Get username")
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
	fmt.Println("HERE")
	if rows.Err() != nil {
		panic(err)
	}

	thread.LikesCount = getLikesFromThreadId(db, thread.Id, true)
	fmt.Println("Get DislikesCount")
	thread.DislikesCount = getLikesFromThreadId(db, thread.Id, false)
	fmt.Println("Get Comments")
	thread.Comments = getComments(db, thread.Id)
	fmt.Println("Get Tags")
	thread.Tags = getTags(db, thread.Id)
	fmt.Println(thread)
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

func EditThreadById(db *sql.DB, title *string, content *string,
	tags []int, threadid int, userId int, token string) error {

	err := checkToken(db, userId, token)

	if err != nil {
		return err
	}

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

func DeleteThread(db *sql.DB, threadId int, token string, userId int) error {
	tx, err := db.Begin()
	if err != nil {
		return errors.New("Unable to begin database transaction")
	}
	defer tx.Rollback()

	err = checkToken(db, userId, token)

	if err != nil {
		return err
	}

	deleteStatement := `
	UPDATE threads
	SET is_deleted = true
	WHERE id = $1
	AND authorid = $2
	AND EXISTS (
		SELECT token
		FROM tokens
		WHERE userid = $2 AND token = $3
	)
	RETURNING id
	`

	rows, err := tx.Query(deleteStatement, threadId, userId, token)

	if err != nil {
		return errors.New("Unable to delete thread")
	}

	//Check if any threads is deleted. throw exception
	//if none is effected.
	isThreadFound := false

	for rows.Next() {
		isThreadFound = true
		break
	}

	if !isThreadFound {
		return errors.New("Thread Not Found")
	}

	// deleteTagsStatement := `
	// DELETE FROM thread_tags
	// WHERE threadid = $1
	// AND EXISTS (
	// 	SELECT id
	// 	FROM threads
	// 	WHERE id = $1
	// 	AND is_deleted = true
	// )
	// `
	// _,  err = tx.Exec(deleteTagsStatement, threadId)

	// if err != nil {
	// 	return err
	// 	// return errors.New("Unable to delete thread")
	// }

	err = tx.Commit()
	if err != nil {
		return errors.New("Unable to commit transaction")
	}

	return nil
}

func PostComment(db *sql.DB, authorid int, content string, parentid int, threadid int) (int, error) {
	newCommentId := getCommentId(db)
	rows, err := db.Query("INSERT INTO Comments (authorid, content, id, is_deleted, parentid, threadid, timestamp) VALUES ($1, $2, $3, $4, $5, $6, $7)", authorid, content, newCommentId, false, parentid, threadid, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return -1, errors.New(err.Error())
	}
	defer rows.Close()
	return newCommentId, nil
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
