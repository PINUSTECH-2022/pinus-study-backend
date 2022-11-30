package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// Loads Database URI from .env file
func loadDbUri() string {
	err := godotenv.Load("database/.env")

	if err != nil {
		panic(err)
	}

	log.Printf("Env value: %s", os.Getenv("DATABASE_URI"))

	return os.Getenv("DATABASE_URI")
}

func GetDb() *sql.DB {
	db, err := sql.Open("postgres", loadDbUri())
	if err != nil {
		panic(err)
	}

	return db
}

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
