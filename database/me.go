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

	username, err := getUsername(db, userid)
	if err != nil {
		return personalInfo, err
	}

	personalInfo.Username = username
	personalInfo.NumberOfQuestionsAsked = getNumberOfThreadsByUser(db, userid) + getNumberOfCommentsByUser(db, userid)
	personalInfo.NumberOfLikesReceived = getNumberOfLikesToUserThreads(db, userid) + getNumberOfLikesToUserComments(db, userid)
	personalInfo.RecentThreads = getRecentThreadsByUser(db, userid)
	personalInfo.MyModules = getModulesSubscribedByUser(db, userid)

	return personalInfo, nil
}
