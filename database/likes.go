package database

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

func GetLikeThread(db *sql.DB, threadid int, userid int) (int, error) {
	rows, err := db.Query("SELECT state FROM Likes_Threads WHERE threadid = $1 AND userid = $2", threadid, userid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	result := 0
	state := false

	for rows.Next() {
		err := rows.Scan(&state)
		if state {
			result = 1
		} else {
			result = -1
		}
		if err != nil {
			panic(err)
		}
	}

	return result, nil
}

// Get the list of user who like a certain thread
func GetListOfLikeThread(db *sql.DB, threadid int) ([]User, error) {
	sql_statement := `
		SELECT u.id, u.username
		FROM likes_threads lt, users u
		WHERE lt.threadId = $1 AND lt.state = TRUE AND lt.userId = u.id;
	`

	rows, err := db.Query(sql_statement, threadid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	like_list := []User{}
	for rows.Next() {
		var newUser User
		rows.Scan(&newUser.Id, &newUser.Username)
		like_list = append(like_list, newUser)
	}

	if rows.Err() != nil {
		panic(rows.Err())
	}

	return like_list, nil
}

// Get the list of user who dislike a certain thread
func GetListOfDislikeThread(db *sql.DB, threadid int) ([]User, error) {
	sql_statement := `
		SELECT u.id, u.username
		FROM likes_threads lt, users u
		WHERE lt.threadId = $1 AND lt.state = FALSE AND lt.userId = u.id;
	`

	rows, err := db.Query(sql_statement, threadid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	dislike_list := []User{}
	for rows.Next() {
		var newUser User
		rows.Scan(&newUser.Id, &newUser.Username)
		dislike_list = append(dislike_list, newUser)
	}

	if rows.Err() != nil {
		panic(rows.Err())
	}

	return dislike_list, nil
}

func SetLikeThread(db *sql.DB, state int, threadid int, userid int) error {
	ResetLikeThread(db, threadid, userid)

	if state == 0 {
		return nil
	}

	var boolState bool
	if state == 1 {
		boolState = true
	} else {
		boolState = false
	}

	rows, err := db.Query("INSERT INTO Likes_Threads (state, threadid, userid) VALUES ($1, $2, $3)", boolState, threadid, userid)
	if err != nil {
		return errors.New(err.Error())
	}
	defer rows.Close()
	return nil
}

func ResetLikeThread(db *sql.DB, threadid int, userid int) error {
	rows, err := db.Query("DELETE FROM Likes_Threads WHERE threadid = $1 AND userid = $2", threadid, userid)
	if err != nil {
		return errors.New(err.Error())
	}
	defer rows.Close()
	return nil
}

func GetLikeComment(db *sql.DB, commentid int, userid int) (int, error) {
	rows, err := db.Query("SELECT state FROM Likes_Comments WHERE commentid = $1 AND userid = $2", commentid, userid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	result := 0
	state := false

	for rows.Next() {
		err := rows.Scan(&state)
		if state {
			result = 1
		} else {
			result = -1
		}
		if err != nil {
			panic(err)
		}
	}

	return result, nil
}

// Get the list of user who like a certain comment
func GetListOfLikeComment(db *sql.DB, commentid int) ([]User, error) {
	sql_statement := `
		SELECT u.id, u.username
		FROM likes_comments lc, users u
		WHERE lc.commentId = $1 AND lc.state = TRUE AND lc.userId = u.id;
	`

	rows, err := db.Query(sql_statement, commentid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	like_list := []User{}
	for rows.Next() {
		var newUser User
		rows.Scan(&newUser.Id, &newUser.Username)
		like_list = append(like_list, newUser)
	}

	if rows.Err() != nil {
		panic(rows.Err())
	}

	return like_list, nil
}

// Get the list of user who dislike a certain comment
func GetListOfDislikeComment(db *sql.DB, commentid int) ([]User, error) {
	sql_statement := `
		SELECT u.id, u.username
		FROM likes_comments lc, users u
		WHERE lc.commentId = $1 AND lc.state = FALSE AND lc.userId = u.id;
	`

	rows, err := db.Query(sql_statement, commentid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	dislike_list := []User{}
	for rows.Next() {
		var newUser User
		rows.Scan(&newUser.Id, &newUser.Username)
		dislike_list = append(dislike_list, newUser)
	}

	if rows.Err() != nil {
		panic(rows.Err())
	}

	return dislike_list, nil
}

func SetLikeComment(db *sql.DB, state int, commentid int, userid int) error {
	ResetLikeComment(db, commentid, userid)

	if state == 0 {
		return nil
	}

	var boolState bool
	if state == 1 {
		boolState = true
	} else {
		boolState = false
	}

	rows, err := db.Query("INSERT INTO Likes_Comments (state, commentid, userid) VALUES ($1, $2, $3)", boolState, commentid, userid)
	if err != nil {
		return errors.New(err.Error())
	}
	defer rows.Close()
	return nil
}

func ResetLikeComment(db *sql.DB, commentid int, userid int) error {
	rows, err := db.Query("DELETE FROM Likes_Comments WHERE commentid = $1 AND userid = $2", commentid, userid)
	if err != nil {
		return errors.New(err.Error())
	}
	defer rows.Close()
	return nil
}
