package database

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

type User struct {
	Id       int
	Username string
}

type UserInfo struct {
	Username               string
	Followers              []User
	Following              []User
	NumberOfFollowers      int
	NumberOfFollowing      int
	NumberOfQuestionsAsked int
	NumberOfLikesReceived  int
	RecentThreads          []Thread
	Modules                []string
}

func GetUserInfoByID(db *sql.DB, userid int) (UserInfo, error) {
	var userInfo UserInfo
	userInfo.Modules = []string{}
	userInfo.RecentThreads = []Thread{}

	sql_statement1 := `
	SELECT u.username, (CASE WHEN s.moduleid IS NULL THEN '' ELSE s.moduleid END) AS module_code
	FROM Users u
	LEFT JOIN Subscribes s ON u.id = s.userid
	WHERE u.id = $1
	`
	rows1, err1 := db.Query(sql_statement1, userid)
	if err1 != nil {
		panic(err1)
	}
	defer rows1.Close()

	for rows1.Next() {
		var module string
		rows1.Scan(&userInfo.Username, &module)
		if module != "" {
			userInfo.Modules = append(userInfo.Modules, module)
		}
	}

	sql_statement2 := `
	SELECT COUNT(lc.userid)
	FROM Users u
	JOIN Comments c ON u.id = c.authorid
	LEFT JOIN Likes_Comments lc ON c.id = lc.commentid AND lc.state = TRUE
	WHERE u.id = $1
	GROUP BY c.id
	`
	rows2, err2 := db.Query(sql_statement2, userid)
	if err2 != nil {
		panic(err2)
	}
	defer rows2.Close()

	for rows2.Next() {
		var comment_like_count int
		rows1.Scan(&comment_like_count)
		userInfo.NumberOfQuestionsAsked += 1
		userInfo.NumberOfLikesReceived += comment_like_count
	}

	sql_statement3 := `
	SELECT t.id, t.title, t.content, t.moduleid, t.authorid, t.timestamp, t.is_deleted,
	COUNT(lt1.userid) AS likes_count, COUNT(lt2.userid) AS dislikes_count
	FROM USERS u
	JOIN Threads t ON u.id = t.authorid
	LEFT JOIN Likes_Threads lt1 ON t.id = lt1.threadid AND lt1.state = TRUE
	LEFT JOIN Likes_Threads lt2 ON t.id = lt2.threadid AND lt2.state = FALSE
	WHERE u.id = $1
	GROUP BY t.id, t.title, t.content, t.moduleid, t.authorid, t.timestamp, t.is_deleted
	ORDER BY t.timestamp, t.id DESC
	`

	rows3, err3 := db.Query(sql_statement3, userid)
	if err3 != nil {
		panic(err3)
	}
	defer rows3.Close()

	var thread_count = 0

	for rows3.Next() {
		var thread Thread

		rows3.Scan(&thread.Id, &thread.Title, &thread.Content, &thread.ModuleId, &thread.AuthorId, &thread.Timestamp,
			&thread.IsDeleted, &thread.LikesCount, &thread.DislikesCount)

		thread_count += 1
		userInfo.NumberOfQuestionsAsked += 1
		userInfo.NumberOfLikesReceived += thread.LikesCount

		// recent threads
		if thread_count <= 5 {
			userInfo.RecentThreads = append(userInfo.RecentThreads, thread)
		}
	}

	followers, err4 := GetFollowers(db, userid)
	if err4 != nil {
		panic(err4)
	}

	following, err5 := GetFollowings(db, userid)
	if err5 != nil {
		panic(err5)
	}

	userInfo.Followers = followers
	userInfo.Following = following
	userInfo.NumberOfFollowers = len(followers)
	userInfo.NumberOfFollowing = len(following)

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

func IsEmailAvailable(db *sql.DB, email string) bool {
	rows, _ := db.Query("SELECT email FROM Users WHERE LOWER(email) = LOWER($1);", email)
	for rows.Next() {
		return false
	}
	return true
}

func IsUsernameAvailable(db *sql.DB, username string) bool {
	rows, _ := db.Query("SELECT username FROM Users WHERE LOWER(username) = LOWER($1)", username)
	for rows.Next() {
		return false
	}
	return true
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
	WHERE LOWER(u.username) = LOWER($1) OR LOWER(u.email) = LOWER($1);
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

func ChangeUsername(db *sql.DB, userid int, newUsername string) error {
	sql_statement := `
	UPDATE Users
	SET username = $1
	WHERE id = $2
	`

	if !IsUsernameAvailable(db, newUsername) {
		return errors.New("Username already exists!")
	}

	_, err := db.Exec(sql_statement, newUsername, userid)
	if err != nil {
		panic(err)
	}

	return nil
}
