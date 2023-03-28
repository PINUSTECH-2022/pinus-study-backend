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

	// Eagerly starts a connection with db
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}