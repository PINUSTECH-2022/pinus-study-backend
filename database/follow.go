package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type Follow struct {
	FollowerId  int
	FollowingId int
	Timestamp   string
}

func FollowUser(db *sql.DB, followerid int, followingid int) error {
	sql_query := `
	INSERT INTO follows (followerid, followingid)
	VALUES ($1, $2)
	`

	_, err := db.Exec(sql_query, followerid, followingid)
	fmt.Println(followerid, followingid)

	if err != nil {
		fmt.Println(err.Error())
		return errors.New("Unable to post follow")
	}

	return nil
}

func UnfollowUser(db *sql.DB, followerid int, followingid int) error {
	sql_query := `
	DELETE FROM follows
	WHERE followerid = $1 AND followingid = $2
	`

	_, err := db.Exec(sql_query, followerid, followingid)
	fmt.Println(followerid, followingid)

	if err != nil {
		fmt.Println(err.Error())
		return errors.New("Unable to unfollow")
	}

	return nil
}
