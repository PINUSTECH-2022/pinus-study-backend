package database

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type UserInfo struct {
	Username               string
	NumberOfQuestionsAsked int
	NumberOfLikesReceived  int
	RecentThreads          []Thread
	Modules              []string
}

func GetUserInfoByID(db *sql.DB, userid int) (UserInfo, error) {
	var userInfo UserInfo

	username, err := getUsername(db, userid)

	if err != nil {
		return userInfo, err
	}

	userInfo.Username = username
	userInfo.NumberOfQuestionsAsked = getNumberOfThreadsByUser(db, userid) + getNumberOfCommentsByUser(db, userid)
	userInfo.NumberOfLikesReceived = getNumberOfLikesToUserThreads(db, userid) + getNumberOfLikesToUserComments(db, userid)
	userInfo.RecentThreads = getRecentThreadsByUser(db, userid)
	userInfo.Modules = getModulesSubscribedByUser(db, userid)

	return userInfo, nil
}

func getUserId(db *sql.DB) int {
	sql_statement := `
	SELECT COUNT(*)
	FROM Users
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

func getUsername(db *sql.DB, userid int) (string, error) {
	sql_statement := `
	SELECT u.username
	FROM Users u
	WHERE u.id = $1
	`
	rows, err := db.Query(sql_statement, userid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var username string
	for rows.Next() {
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
	}

	if rows.Err() != nil {
		panic(err)
	}

	if username == "" {
		return "", errors.New("User not found")
	}

	return username, nil
}

func getUserIdFromNameOrEmail(db *sql.DB, nameOrEmail string) (int, error) {
	sql_statement := `
	SELECT u.id
	FROM Users u
	WHERE u.username = $1 OR u.email = $1
	`
	rows, err := db.Query(sql_statement, nameOrEmail)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	userid := -1
	for rows.Next() {
		err := rows.Scan(&userid)
		if err != nil {
			panic(err)
		}
	}

	if rows.Err() != nil {
		panic(err)
	}

	if userid == -1 {
		return -1, errors.New("Userid not found")
	}

	return userid, nil
}
