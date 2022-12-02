package database

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

func GetSubscribers(db *sql.DB, moduleid string) []int {
	rows, err := db.Query("SELECT userid FROM Subscribes WHERE moduleid = $1", moduleid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var subscribers []int

	for rows.Next() {
		var subscriber int
		err := rows.Scan(&subscriber)
		if err != nil {
			panic(err)
		}
		subscribers = append(subscribers, subscriber)
	}

	if rows.Err() != nil {
		panic(err)
	}

	return subscribers
}

func DoesSubscribe(db *sql.DB, moduleid string, userid int) (bool, error) {
	rows, err := db.Query("SELECT userid FROM Subscribes WHERE moduleid = $1 AND userid = $2", moduleid, userid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	result := -1

	for rows.Next() {
		err := rows.Scan(&result)
		if err != nil {
			panic(err)
		}
	}

	if rows.Err() != nil {
		panic(err)
	}

	return result != -1, nil
}

func Subscribe(db *sql.DB, moduleid string, userid int) error {
	rows, err := db.Query("INSERT INTO Subscribes (moduleid, userid) VALUES ($1, $2)", moduleid, userid)
	if err != nil {
		return errors.New(err.Error())
	}
	defer rows.Close()
	return nil
}

func Unsubscribe(db *sql.DB, moduleid string, userid int) error {
	rows, err := db.Query("DELETE FROM Subscribes WHERE moduleid = $1 AND userid = $2", moduleid, userid)
	if err != nil {
		return errors.New(err.Error())
	}
	defer rows.Close()
	return nil
}
