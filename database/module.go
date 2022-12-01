package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Module struct {
	Id   string
	Name string
	Desc string
}

func GetModules(db *sql.DB) []Module {
	rows, err := db.Query("SELECT * FROM Modules")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var modules []Module

	for rows.Next() {
		var mod Module
		err := rows.Scan(&mod.Id, &mod.Name, &mod.Desc)
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
