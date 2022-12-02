package database

import (
	"database/sql"

	_ "github.com/lib/pq"
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

func GetCommentById(db *sql.DB, id int) Comment {
	sql_statement := `
	SELECT c.content, u.id, u.username, c.is_deleted, c.timestamp
	FROM Comments c JOIN Users u ON c.authorid = u.id
	WHERE c.id = $1
	`
	rows, err := db.Query(sql_statement, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var comment Comment

	for rows.Next() {
		err := rows.Scan(&comment.Content, &comment.AuthorId,
			&comment.Username, &comment.IsDeleted, &comment.Timestamp)
		if err != nil {
			panic(err)
		}
	}

	if rows.Err() != nil {
		panic(err)
	}

	comment.Likes = getLikesFromCommentId(db, id, true)
	comment.Dislikes = getLikesFromCommentId(db, id, false)
	comment.CommentChilds = getChildrensFromCommentId(db, id)

	return comment
}

// return true if it runs correctly
func DeleteCommentById(db *sql.DB, commentId int, userId int,
	token string) bool {
	if !isAuthorized(db, commentId, userId, token) {
		return false
	}

	sql_statement := `
	UPDATE Comments
	SET is_deleted = TRUE
	WHERE id = $1
	`
	_, err := db.Exec(sql_statement, commentId)

	return err == nil
}

// return true if it runs correctly
func UpdateCommentById(db *sql.DB, commentId int, content string,
	userId int, token string) bool {
	if !isAuthorized(db, commentId, userId, token) {
		return false
	}

	if content == "" {
		return false
	}

	sql_statement := `
	UPDATE Comments
	SET content = $2
	WHERE id = $1
		`
	_, err := db.Exec(sql_statement, commentId, content)
	if err != nil {
		panic(err)
	}

	return err == nil
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

// might want to change it later
func isAuthorized(db *sql.DB, commentId int, userId int, token string) bool {
	var status bool
	sql_statement := `
	SELECT 1
	FROM Tokens t JOIN Comments c ON t.userid = c.authorid
	WHERE c.id = $1 AND c.authorid = $2 AND t.token = $3
	`
	rows, err := db.Query(sql_statement, commentId, userId, token)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		status = true
	}

	return status
}
