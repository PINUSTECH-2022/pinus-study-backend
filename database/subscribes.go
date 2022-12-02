package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func GetSubcsribers(db *sql.DB, moduleid string) []int {
	rows, err := db.Query("SELECT userid FROM Subscribes WHERE moduleid = $1", moduleid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var subcsribers []int

	for rows.Next() {
		var subcsriber int
		err := rows.Scan(&subcsriber)
		if err != nil {
			panic(err)
		}
		subcsribers = append(subcsribers, subcsriber)
	}

	if rows.Err() != nil {
		panic(err)
	}

	return subcsribers
}

func Subcsribe(db *sql.DB, moduleid string, userid int) {
	rows, err := db.Query("INSERT INTO Subscribes (moduleid, userid) VALUES ($1, $2)", moduleid, userid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
}

func Unsubcsribe(db *sql.DB, moduleid string, userid int) {
	rows, err := db.Query("DELETE FROM Subscribes WHERE moduleid = $1 AND userid = $2", moduleid, userid)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
}
