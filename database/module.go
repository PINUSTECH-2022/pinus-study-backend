package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Module struct {
	Id              string
	Name            string
	Desc            string
	SubscriberCount int
	Threads         []Thread
	ReviewCount     int
}

func GetModules(db *sql.DB, keyword string, page int) []Module {
	fmt.Println("Executing SQL. Keyword: ", keyword)
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

	rows, err := db.Query(sql_statement, keyword, 12*(page-1))
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
	fmt.Println(modules)
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

func getReviewCount(db *sql.DB, moduleid string) int {
	rows, err := db.Query(`SELECT COUNT(*)
		FROM Reviews AS R
		WHERE R.moduleId = $1 AND R.is_deleted = false`,
		moduleid)

	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			panic(err)
		}
	}

	return count
}

func GetModuleByModuleId(db *sql.DB, moduleid string) Module {
	query := fmt.Sprintf(`
	SELECT M.id, M.name, M.description, COUNT(S.moduleid)
	FROM Modules AS M 
	LEFT JOIN Threads AS T ON M.id = T.moduleid 
	LEFT JOIN Subscribes AS S ON S.moduleid = M.id 
	WHERE M.id = '%s' 
	GROUP BY M.id, M.name, M.description, T.id, T.title, T.content, T.moduleid, T.authorid, T.timestamp, T.is_deleted;
	`, moduleid)

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer rows.Close()

	var mod Module
	for rows.Next() {
		err := rows.Scan(&mod.Id, &mod.Name, &mod.Desc, &mod.SubscriberCount)
		if err != nil {
			panic(err)
		}
	}

	query = fmt.Sprintf(`
	SELECT T.id, T.title, T.content, T.moduleid, T.authorid, U.username, T.timestamp, T.is_deleted, T.likes_count, T.dislikes_count, T.comments_count
	FROM Modules AS M 
	LEFT JOIN Threads AS T ON M.id = T.moduleid 
	LEFT JOIN Subscribes AS S ON S.moduleid = M.id 
	LEFT JOIN Users AS U ON U.id = T.authorid
	WHERE M.id = '%s';
	`, moduleid)

	rows, err = db.Query(query)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var thread Thread
		err := rows.Scan(&thread.Id, &thread.Title, &thread.Content, &thread.ModuleId, &thread.AuthorId, &thread.Username,
			&thread.Timestamp, &thread.IsDeleted, &thread.LikesCount, &thread.DislikesCount, &thread.CommentsCount)
		if err != nil {
			break
		}
		mod.Threads = append(mod.Threads, thread)
	}

	mod.ReviewCount = getReviewCount(db, moduleid)

	fmt.Println(mod)
	return mod
}

func PostThread(db *sql.DB, authorid int, content string, title string, tags []int, moduleid string) (int, error) {
	fmt.Println("Posting thread...")
	fmt.Println(authorid, content, title, tags, moduleid)
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("Error in initializing db: ", err.Error())
		return -1, errors.New("Unable to begin database transaction")
	}
	defer tx.Rollback()

	newThreadID := getThreadId(tx)

	_, err = tx.Exec("INSERT INTO Threads (authorid, content, id, moduleid, title) VALUES ($1, $2, $3, $4, $5)", authorid, content, newThreadID, strings.ToUpper(moduleid), title)
	if err != nil {
		fmt.Println("Error in inserting thread into db: ", err.Error())
		return -1, errors.New("Thread data is malformed.")
	}

	for _, tagId := range tags {
		_, err := tx.Exec("INSERT INTO Thread_Tags VALUES ($1, $2)", newThreadID, tagId)
		if err != nil {
			fmt.Println("Error in inserting thread tag into db: ", err.Error())
			return -1, errors.New("Tag data is malformed.")
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Error in commiting thread posted: ", err.Error())
		return -1, errors.New("Unable to commit transaction")
	}
	fmt.Println("Posted...")
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
