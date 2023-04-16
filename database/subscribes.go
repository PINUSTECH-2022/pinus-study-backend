package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type Subscriber struct {
	Id       int
	Username string
}

func GetSubscribers(db *sql.DB, moduleid string) []Subscriber {
	rows, err := db.Query(fmt.Sprintf("SELECT S.userid, U.username FROM Subscribes AS S, Users AS U WHERE S.moduleid = '%s' AND S.userid = U.id;", moduleid))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var subscribers []Subscriber

	for rows.Next() {
		var subscriber Subscriber
		err := rows.Scan(&subscriber.Id, &subscriber.Username)
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

func getModulesSubscribedByUser(db *sql.DB, userid int) []string {
	sql_statement := `
	SELECT s.moduleid
	FROM Subscribes s
	WHERE s.userid = $1
	`
	rows, err := db.Query(sql_statement, userid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	subscribedModules := []string{}
	for rows.Next() {
		var mod string
		err := rows.Scan(&mod)
		if err != nil {
			panic(err)
		}
		subscribedModules = append(subscribedModules, mod)
	}

	if rows.Err() != nil {
		panic(err)
	}

	return subscribedModules
}
