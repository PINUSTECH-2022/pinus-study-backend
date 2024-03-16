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

func GetLikeReview(db *sql.DB, reviewid int, userid int) (int, error) {
	rows, err := db.Query("SELECT state FROM Likes_reviews WHERE reviewid = $1 AND userid = $2", reviewid, userid)
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

func SetLikeReview(db *sql.DB, state int, reviewid int, userid int) error {
	ResetLikeThread(db, reviewid, userid)

	if state == 0 {
		return nil
	}

	var boolState bool
	if state == 1 {
		boolState = true
	} else {
		boolState = false
	}

	rows, err := db.Query("INSERT INTO Likes_Reviews (state, reviewid, userid) VALUES ($1, $2, $3)", boolState, reviewid, userid)
	if err != nil {
		return errors.New(err.Error())
	}
	defer rows.Close()
	return nil
}

func ResetLikeReview(db *sql.DB, reviewid int, userid int) error {
	rows, err := db.Query("DELETE FROM Likes_Reviews WHERE reviewid = $1 AND userid = $2", reviewid, userid)
	if err != nil {
		return errors.New(err.Error())
	}
	defer rows.Close()
	return nil
}
