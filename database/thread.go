package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"

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
	CommentsCount int
	IsDeleted     bool
	Comments      []int
	Tags          []int
}

func GetThreadById(db *sql.DB, threadid string) Thread {
	threadidInt, err := strconv.Atoi(threadid)
	if err != nil {
		panic(err)
	}

	query := fmt.Sprintf(`
	SELECT T.id, T.title, T.content, T.moduleid, T.authorid, T.timestamp, T.is_deleted,  T.likes_count, T.dislikes_count, T.comments_count, U.username, C.id 
	FROM Threads AS T 
	JOIN Users AS U ON T.authorid = U.id 
	LEFT JOIN Comments AS C ON C.threadid = T.id 
	WHERE T.id = %d AND (C.parentid = 0 OR C.parentid IS NULL);`,
		threadidInt)

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer rows.Close()

	var thread Thread
	for rows.Next() {
		var commentId sql.NullInt64
		err := rows.Scan(&thread.Id, &thread.Title, &thread.Content, &thread.ModuleId, &thread.AuthorId, &thread.Timestamp,
			&thread.IsDeleted, &thread.LikesCount, &thread.DislikesCount, &thread.CommentsCount, &thread.Username, &commentId)
		if err != nil {
			panic(err)
		}
		var commentIdNotNull int
		if commentId.Valid {
			commentIdNotNull = int(commentId.Int64)
			thread.Comments = append(thread.Comments, commentIdNotNull)
		}
	}
	return thread

	// var wg sync.WaitGroup

	// thread_c := make(chan Thread, 1)
	// likeCount_c := make(chan int, 1)
	// dislikeCount_c := make(chan int, 1)
	// comments_c := make(chan []int, 1)
	// tags_c := make(chan []int, 1)

	// wg.Add(5)
	// go func(db *sql.DB, threadid int, c chan Thread) {
	// 	defer wg.Done()
	// 	defer close(c)
	// 	var thread Thread

	// 	err = db.QueryRow(
	// 	`SELECT t.id, t.title, t.content, t.moduleid,
	// 	t.authorid, t.timestamp, t.is_deleted, u.username
	// 	FROM Threads as t, Users as u
	// 	WHERE u.id = t.id AND t.id = $1`,
	// 	threadidInt).Scan(&thread.Id, &thread.Title, &thread.Content, &thread.ModuleId, &thread.AuthorId, &thread.Timestamp, &thread.IsDeleted, &thread.Username)

	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	c <- thread
	// } (db, threadidInt, thread_c)
	// go func(db *sql.DB, threadid int, c chan int) {
	// 	defer wg.Done()
	// 	defer close(c)
	// 	c <- getLikesFromThreadId(db, threadid, true)
	// } (db, threadidInt, likeCount_c)
	// go func(db *sql.DB, threadid int, c chan int) {
	// 	defer wg.Done()
	// 	defer close(c)
	// 	c <- getLikesFromThreadId(db, threadid, false)
	// } (db, threadidInt, dislikeCount_c)
	// go func(db *sql.DB, threadid int, c chan []int) {
	// 	defer wg.Done()
	// 	defer close(c)
	// 	c <- getComments(db, threadid)
	// } (db, threadidInt, comments_c)
	// go func(db *sql.DB, threadid int, c chan []int) {
	// 	defer wg.Done()
	// 	defer close(c)
	// 	c <- getTags(db, threadid)
	// } (db, threadidInt, tags_c)

	// wg.Wait()

	// thread := <-thread_c
	// thread.LikesCount = <-likeCount_c
	// thread.DislikesCount = <-dislikeCount_c
	// thread.Comments = <-comments_c
	// thread.Tags = <-tags_c
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
	pid := sql.NullInt64{
		Int64: int64(parentid),
		Valid: parentid != 0,
	}
	var newCommentId int
	err := db.QueryRow(`INSERT INTO 
	Comments (authorid, content, id, is_deleted, parentid, threadid) 
	VALUES ($1, $2, (SELECT COUNT(*)
	FROM Comments) + 1, $3, $4, $5) RETURNING id`, authorid, content, false, pid, threadid).Scan(&newCommentId)
	if err != nil {
		fmt.Println(err.Error())
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
