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
	COUNT(lt1.userid) AS likes_count, COUNT(lt2.userid) AS dislikes_count, 
	(CASE WHEN c.id IS NULL THEN -1 ELSE c.id END) AS comment_id, 
	(CASE WHEN tt.tagId IS NULL THEN -1 ELSE tt.tagId END) AS tag_id
	FROM USERS u
	JOIN Threads t ON u.id = t.authorid
	LEFT JOIN Comments c ON t.id = c.authorid
	LEFT JOIN Likes_Threads lt1 ON t.id = lt1.threadid AND lt1.state = FALSE
	LEFT JOIN Likes_Threads lt2 ON t.id = lt2.threadid AND lt2.state = FALSE
	LEFT JOIN Thread_Tags tt ON t.id = tt.threadid
	WHERE u.id = $1
	GROUP BY t.id, t.title, t.content, t.moduleid, t.authorid, t.timestamp, t.is_deleted,
	c.id, tt.tagId
	ORDER BY t.timestamp, t.id DESC
	`

	rows3, err3 := db.Query(sql_statement3, userid)
	if err3 != nil {
		panic(err3)
	}
	defer rows3.Close()

	var prev_thread_id = -1
	var thread_count = 0
	comment_list := make(map[int]int)
	tag_list := make(map[int]int)

	for rows3.Next() {
		var thread Thread
		var comment_id int
		var tag_id int

		rows3.Scan(&thread.Id, &thread.Title, &thread.Content, &thread.ModuleId, &thread.AuthorId, &thread.Timestamp,
			&thread.IsDeleted, &thread.LikesCount, &thread.DislikesCount, &comment_id, &tag_id)
		thread.Comments = []int{}
		thread.Tags = []int{}

		if prev_thread_id != thread.Id {
			thread_count += 1
			prev_thread_id = thread.Id
			personalInfo.NumberOfQuestionsAsked += 1
			personalInfo.NumberOfLikesReceived += thread.LikesCount

			if thread_count <= 5 {
				// reset the map
				comment_list = make(map[int]int)
				tag_list = make(map[int]int)

				personalInfo.RecentThreads = append(personalInfo.RecentThreads, thread)
			}
		}
		if thread_count <= 5 {
			if comment_id != -1 && comment_list[comment_id] == 0 {
				personalInfo.RecentThreads[len(personalInfo.RecentThreads)-1].Comments = append(
					personalInfo.RecentThreads[len(personalInfo.RecentThreads)-1].Comments, comment_id)
				comment_list[comment_id] = 1
			}
			if tag_id != -1 && tag_list[tag_id] == 0 {
				personalInfo.RecentThreads[len(personalInfo.RecentThreads)-1].Tags = append(
					personalInfo.RecentThreads[len(personalInfo.RecentThreads)-1].Tags, tag_id)
				tag_list[tag_id] = 1
			}
		}
	}

	return personalInfo, nil
}
