package database

import (
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

type Module struct {
	Id              string
	Name            string
	Desc            string
	SubscriberCount int
}

func GetModules(db *sql.DB, keyword string, page int) []Module {
	sql_statement := `
	SELECT * FROM MODULES
	WHERE id LIKE '%' || UPPER($1) || '%'
	OFFSET $2
	LIMIT 10
	`

	rows, err := db.Query(sql_statement, keyword, 10*(page-1))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var modules []Module

	for rows.Next() {
		var mod Module
		err := rows.Scan(&mod.Id, &mod.Name, &mod.Desc)
		mod.SubscriberCount = getSubscriberCount(db, mod.Id)
		if err != nil {
			panic(err)
		}
		modules = append(modules, mod)
	}

	if rows.Err() != nil {
		panic(err)
	}

	return modules
}

func getSubscriberCount(db *sql.DB, moduleid string) int {
	sql_statement := `
	SELECT COUNT(*)
	FROM Subscribes
	WHERE moduleid = $1
	`
	rows, err := db.Query(sql_statement, moduleid)
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

	return count
}

func GetModuleByModuleId(db *sql.DB, moduleid string) Module {
	rows, err := db.Query("SELECT * FROM Modules WHERE id = $1", moduleid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var mod Module

	for rows.Next() {
		err := rows.Scan(&mod.Id, &mod.Name, &mod.Desc)
		mod.SubscriberCount = getSubscriberCount(db, mod.Id)
		if err != nil {
			panic(err)
		}
	}

	if rows.Err() != nil {
		panic(err)
	}

	return mod
}

func PostThread(db *sql.DB, authorid int, content string, title string, tags []int, moduleid string) error {

	tx, err := db.Begin()
	if (err != nil) {
		return errors.New("Unable to begin database transaction")
	}
	defer tx.Rollback()

	newThreadID := getThreadId(tx)

	_, err = tx.Exec("INSERT INTO Threads (authorid, content, id, moduleid, timestamp, title) VALUES ($1, $2, $3, $4, $5, $6)", authorid, content, newThreadID, moduleid, time.Now().Format("2006-01-02 15:04:05"), title)
	if err != nil {
		return errors.New("Thread data is malformed.")
	}

	for _, tagId := range tags {
		_, err := tx.Exec("INSERT INTO Thread_Tags VALUES ($1, $2)", newThreadID, tagId)
		if err != nil {
			return errors.New("Tag data is malformed.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.New("Unable to commit transaction")
	}

	return nil
}

func getThreadId(tx *sql.Tx) int {
	sql_statement := `
	SELECT COUNT(*)
	FROM Threads
	`
	rows, err := tx.Query(sql_statement)
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
