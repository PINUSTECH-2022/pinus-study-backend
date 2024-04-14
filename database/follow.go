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
		return errors.New("unable to post follow")
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
		return errors.New("unable to unfollow")
	}

	return nil
}

func GetFollowers(db *sql.DB, userid int) ([]User, error) {
	sql_statement := `
	SELECT f.followerid, u.username
	FROM follows f, users u
	WHERE f.followerid = u.id AND f.followingid = $1
	`
	rows, err := db.Query(sql_statement, userid)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("something went wrong when getting followers")
	}
	defer rows.Close()

	var followers []User
	for rows.Next() {
		var follower User
		err := rows.Scan(&follower.Id, &follower.Username)
		if err != nil {
			fmt.Println(err.Error())
			return nil, errors.New("something went wrong when getting followers")
		}
		followers = append(followers, follower)
	}

	if rows.Err() != nil {
		fmt.Println(err.Error())
		return nil, errors.New("something went wrong when getting followers")
	}

	return followers, nil
}

func GetFollowings(db *sql.DB, userid int) ([]User, error) {
	sql_statement := `
	SELECT f.followingid, u.username
	FROM follows f, users u
	WHERE f.followingid = u.id AND f.followerid = $1
	`
	rows, err := db.Query(sql_statement, userid)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("something went wrong when getting followings")
	}
	defer rows.Close()

	var followings []User
	for rows.Next() {
		var following User
		err := rows.Scan(&following.Id, &following.Username)
		if err != nil {
			fmt.Println(err.Error())
			return nil, errors.New("something went wrong when getting followings")
		}
		followings = append(followings, following)
	}

	if rows.Err() != nil {
		fmt.Println(err.Error())
		return nil, errors.New("something went wrong when getting followings")
	}

	return followings, nil
}

// Get a list of thread id which is posted by the user's following
func GetFollowingsThreads(db *sql.DB, userid int) ([]int, error) {
	sql_statement := `
	SELECT t.id
	FROM threads t, follows f
	WHERE f.followerid = $1 AND f.followingid = t.authorid AND NOT t.is_deleted;
	`

	rows, err := db.Query(sql_statement, userid)
	if err != nil {
		fmt.Println(err.Error())
		return nil, errors.New("something went wrong when getting followings' threads")
	}
	defer rows.Close()

	threads := []int{}
	for rows.Next() {
		var threadid int
		err1 := rows.Scan(&threadid)
		if err1 != nil {
			fmt.Println(err.Error())
			return nil, errors.New("something went wrong when getting followings' threads")
		}
		threads = append(threads, threadid)
	}

	if rows.Err() != nil {
		fmt.Println(err.Error())
		return nil, errors.New("something went wrong when getting followings' threads")
	}

	return threads, nil
}
