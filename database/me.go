package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PersonalInfo struct {
	Username               string
	NumberOfQuestionsAsked int
	NumberOfLikesReceived  int
	RecentThreads          []Thread
	MyModules              []string
}

func GetPersonalInfo(db *sql.DB, userid int) (PersonalInfo, error) {
	var personalInfo PersonalInfo
	personalInfo.MyModules = []string{}
	personalInfo.RecentThreads = []Thread{}

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
		rows1.Scan(&personalInfo.Username, &module)
		if module != "" {
			personalInfo.MyModules = append(personalInfo.MyModules, module)
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
		personalInfo.NumberOfQuestionsAsked += 1
		personalInfo.NumberOfLikesReceived += comment_like_count
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
		personalInfo.NumberOfQuestionsAsked += 1
		personalInfo.NumberOfLikesReceived += thread.LikesCount

		// recent threads
		if thread_count <= 5 {
			personalInfo.RecentThreads = append(personalInfo.RecentThreads, thread)
		}
	}

	return personalInfo, nil
}
