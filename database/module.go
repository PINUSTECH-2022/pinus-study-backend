package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type Module struct {
	Id              string
	Name            string
	Desc            string
	SubscriberCount int
	Threads         []Thread
}

func GetModules(db *sql.DB, keyword string, page int) []Module {
	fmt.Println("Executing ")
	sql_statement := `
	SELECT Modules.id, COUNT(Threads.moduleid) AS popularity
	FROM Modules
	LEFT JOIN Threads ON Modules.id = Threads.moduleid
	WHERE Modules.id LIKE '%' || UPPER($1) || '%'
	GROUP BY Modules.id
	ORDER BY popularity DESC
	OFFSET $2
	LIMIT 12
	`

	rows, err := db.Query(sql_statement, keyword, 10*(page-1))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var modules []Module

	for rows.Next() {
		var mod Module
		var popularity int
		err := rows.Scan(&mod.Id, &popularity)
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
	fmt.Println("a")
	rows, err := db.Query("SELECT * FROM Modules WHERE id = $1", moduleid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var mod Module
	fmt.Println("b")
	for rows.Next() {
		err := rows.Scan(&mod.Id, &mod.Name, &mod.Desc)
		mod.SubscriberCount = getSubscriberCount(db, mod.Id)
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(mod)
	if rows.Err() != nil {
		panic(err)
	}
	fmt.Println(mod.Id)
	thread_ids, err := db.Query("SELECT id FROM Threads WHERE moduleid = $1", mod.Id)
	if err != nil {
		panic(err)
	}
	defer thread_ids.Close()
	fmt.Println("TEST")
	for thread_ids.Next() {
		var thread_id int
		err := thread_ids.Scan(&thread_id)
		if err != nil {
			panic(err)
		}
		fmt.Println(thread_id)
		mod.Threads = append(mod.Threads, GetThreadById(db, strconv.Itoa(thread_id)))
	}

	fmt.Println(mod)
	return mod
}

func PostThread(db *sql.DB, authorid int, content string, title string, tags []int, moduleid string) (int, error) {

	tx, err := db.Begin()
	if err != nil {
		return -1, errors.New("Unable to begin database transaction")
	}
	defer tx.Rollback()

	newThreadID := getThreadId(tx)

	_, err = tx.Exec("INSERT INTO Threads (authorid, content, id, moduleid, timestamp, title) VALUES ($1, $2, $3, $4, $5, $6)", authorid, content, newThreadID, moduleid, time.Now().Format("2006-01-02 15:04:05"), title)
	if err != nil {
		return -1, errors.New("Thread data is malformed.")
	}

	for _, tagId := range tags {
		_, err := tx.Exec("INSERT INTO Thread_Tags VALUES ($1, $2)", newThreadID, tagId)
		if err != nil {
			return -1, errors.New("Tag data is malformed.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return -1, errors.New("Unable to commit transaction")
	}

	return newThreadID, nil
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
