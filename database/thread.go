package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"sync"
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
	threadidInt, err := strconv.Atoi(threadid)
	if err != nil {
		panic(err)
	}

	// query := `SELECT T.id, T.title, T.content, T.moduleid, T.authorid, T.timestamp, T.is_deleted, 
	//	(SELECT COUNT(*) FROM Likes_threads AS LT WHERE LT.threadid = T.id AND LT.state = true) AS likes_count, 
	//	(SELECT COUNT(*) FROM Likes_threads AS LT WHERE LT.threadid = T.id AND LT.state = false) AS dislikes_count, U.username, C.id 
	//	FROM Threads AS T 
	//	JOIN Users AS U ON T.authorid = U.id 
	//	LEFT JOIN Comments AS C ON C.threadid = T.id 
	//	WHERE T.id = $1;`

	var thread Thread
	var wg sync.WaitGroup

	comments_c := make(chan []int, 1)

	err = db.QueryRow(
	`SELECT t.id, t.title, t.content, t.moduleid, 
	t.authorid, t.timestamp, t.is_deleted, u.username,
	(SELECT COUNT(*) FROM Likes_Threads WHERE state=TRUE AND threadid = t.id),
	(SELECT COUNT(*) FROM Likes_Threads WHERE state=FALSE AND threadid = t.id)
	FROM Threads as t, Users as u
	WHERE u.id = t.authorid AND t.id = $1`, 
	threadidInt).Scan(&thread.Id, &thread.Title, &thread.Content, &thread.ModuleId, 
		&thread.AuthorId, &thread.Timestamp, &thread.IsDeleted, &thread.Username,
	&thread.LikesCount, &thread.DislikesCount)

	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go func(db *sql.DB, threadid int, c chan []int) {
		defer wg.Done()
		defer close(c)
		c <- getComments(db, threadid)
	} (db, threadidInt, comments_c)

	thread.Tags = getTags(db, thread.Id)
	
	wg.Wait()
  
	thread.Comments = <-comments_c

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

	var count int
	err := db.QueryRow(sql_statement, userid).Scan(&count)
	if err != nil {
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
	var count int

	err := db.QueryRow(sql_statement, status, id).Scan(&count)
	if err != nil {
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

	tx, err := db.Begin()
	if err != nil {
		return errors.New("Unable to begin database transaction")
	}
	defer tx.Rollback()

	// Bad authentication method, but works.
	var modifiedIns int
	err = tx.QueryRow(`UPDATE Threads 
	SET title = COALESCE($1, title), content = COALESCE($2, content) 
	WHERE id = $3 AND authorid = $4 AND EXISTS (
		SELECT token
		FROM tokens
		WHERE userid = $4 AND token = $5) RETURNING id`, title, content, threadid, userId, token).Scan(&modifiedIns)

	if err != nil {
		return errors.New("Title or content has improper formatting, or you do not have the necessarily priviliges to modify this asset")
	}

	if tags != nil {
		_, err := tx.Exec("DELETE FROM Thread_Tags WHERE threadid = $1", threadid)
		if err != nil {
			return errors.New("Unable to delete thread tags")
		}

		if len(tags) > 0 {
			sql_statement := "INSERT INTO Thread_Tags VALUES "
			vals := []any{}
			vals = append(vals, threadid)

			for i, tagId := range tags {
				sql_statement += fmt.Sprintf("($1, $%d),", i+2)
				vals = append(vals, tagId)
			}

			_, err = tx.Exec(sql_statement[0:(len(sql_statement)-1)], vals...)
			if err != nil {
				fmt.Println(err)
				return errors.New("Tags have improper formatting")
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

	/*err := checkToken(db, userId, token)
	if err != nil {
		return err
	}*/

	// Terrible way of doing auth. Should make a middleware that automates token checking.
	deleteStatement := `
	UPDATE threads
	SET is_deleted = true
	WHERE id = $1
	AND authorid = $2
	AND is_deleted = false
	AND EXISTS (
		SELECT token
		FROM tokens
		WHERE userid = $2 AND token = $3
	)
	RETURNING id
	`

	var tmpId int
	err := db.QueryRow(deleteStatement, threadId, userId, token).Scan(&tmpId)

	if err == sql.ErrNoRows {
		return errors.New("Thread is not found or has been deleted")
	} else if err != nil {
		return errors.New("Unable to delete thread")
	}

	return nil
}

func PostComment(db *sql.DB, authorid int, content string, parentid int, threadid int) (int, error) {
	var newCommentId int
	err := db.QueryRow(`INSERT INTO 
	Comments (authorid, content, id, is_deleted, parentid, threadid, timestamp) 
	VALUES ($1, $2, (SELECT COUNT(*)
	FROM Comments) + 1, $3, $4, $5, $6) RETURNING id`, authorid, content, false, parentid, threadid, time.Now().Format("2006-01-02 15:04:05")).Scan(&newCommentId)
	if err != nil {
		return -1, errors.New("Unable to post comment")
	}
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
