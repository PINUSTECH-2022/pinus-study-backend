package database

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Comment struct {
	Content       string
	AuthorId      int
	Username      string
	Likes         int
	Dislikes      int
	IsDeleted     bool
	Timestamp     string
	CommentChilds []int
}

// Get comment by it's id
func GetCommentById(db *sql.DB, id int) Comment {
	sql_statement := `
	SELECT c.content, u.id, u.username, c.is_deleted, c.timestamp, c.likes_count, c.dislikes_count,
		(
			SELECT ARRAY(
			 	SELECT DISTINCT ch.id
				FROM Comments ch
				WHERE ch.parentid = c.id
			)
		) as childComments
	FROM Comments c JOIN Users u ON c.authorid = u.id
	WHERE c.id = $1
	`
	rows, err := db.Query(sql_statement, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var comment Comment
	var childComments []int64

	for rows.Next() {
		err := rows.Scan(&comment.Content, &comment.AuthorId,
			&comment.Username, &comment.IsDeleted, &comment.Timestamp,
			&comment.Likes, &comment.Dislikes, pq.Array(&childComments))
		if err != nil {
			panic(err)
		}
	}

	if rows.Err() != nil {
		panic(err)
	}

	comment.CommentChilds = make([]int, len(childComments))

	for i, id := range childComments {
		comment.CommentChilds[i] = int(id)
	}

	return comment
}

func getNumberOfCommentsByUser(db *sql.DB, userid int) int {
	sql_statement := `
	SELECT COUNT(*)
	FROM Comments c
	WHERE c.authorid = $1
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

func DeleteCommentById(db *sql.DB, commentId int, userId int) error {
	if !isAuthorized(db, commentId, userId) {
		return errors.New("Not authorized")
	}

	sql_statement := `
	UPDATE Comments
	SET is_deleted = TRUE
	WHERE id = $1
	`
	_, err := db.Exec(sql_statement, commentId)

	return err
}

func UpdateCommentById(db *sql.DB, commentId int, userId int, content string) error {
	if !isAuthorized(db, commentId, userId) {
		return errors.New("Not authorized")
	}

	if content == "" {
		return errors.New("Comment must not be empty")
	}

	sql_statement := `
	UPDATE Comments
	SET content = $2
	WHERE id = $1
		`
	_, err := db.Exec(sql_statement, commentId, content)

	return err
}

// if status true, return number of likes
// else if status is false, return number of dislikes
func getLikesFromCommentId(db *sql.DB, id int, status bool) int {
	sql_statement := `
	SELECT COUNT(*)
	FROM Likes_Comments
	WHERE state = $1 AND commentid = $2
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

func getNumberOfLikesToUserComments(db *sql.DB, userid int) int {
	sql_statement := `
	SELECT c.id
	FROM Comments c
	WHERE c.authorid = $1
	`
	rows, err := db.Query(sql_statement, userid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var commentid int
		rows.Scan(&commentid)
		count += getLikesFromThreadId(db, commentid, true)
	}

	if rows.Err() != nil {
		panic(err)
	}

	return count
}

func getChildrensFromCommentId(db *sql.DB, id int) []int {
	sql_statement := `
	SELECT c.id
	FROM Comments c
	WHERE c.parentid = $1
	`
	rows, err := db.Query(sql_statement, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var childs []int
	for rows.Next() {
		var child int
		err := rows.Scan(&child)
		if err != nil {
			panic(err)
		}
		childs = append(childs, child)
	}

	if rows.Err() != nil {
		panic(err)
	}

	return childs
}

func isAuthorized(db *sql.DB, commentId int, userId int) bool {
	var status bool
	sql_statement := `
	SELECT 1 FROM
	Comments c JOIN Users u
	ON c.authorid = u.id
	WHERE c.id = $1 AND u.id = $2
	`
	rows, err := db.Query(sql_statement, commentId, userId)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		status = true
	}

	return status
}
