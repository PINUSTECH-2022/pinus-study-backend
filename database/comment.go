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

func GetCommentById(db *sql.DB, id string) Comment {
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

	if rows == nil {
		return comment
	}

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

// if status true, return number of likes
// else if status is false, return number of dislikes
func getLikesFromCommentId(db *sql.DB, id string, status bool) int {
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

func getChildrensFromCommentId(db *sql.DB, id string) []int {
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
