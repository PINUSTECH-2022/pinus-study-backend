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

	dbUri := os.Getenv("DATABASE_URI")

	if dbUri != "" {
		return dbUri
	}

	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}

	dbUri = os.Getenv("DATABASE_URI")
	log.Printf("Env value: %s", dbUri)

	return dbUri
}

func GetDb() *sql.DB {
	db, err := sql.Open("postgres", loadDbUri())
	if err != nil {
		panic(err)
	}

	// Increase maximum idle connections to improve latency.
	db.SetMaxIdleConns(5)

	// Eagerly starts a connection with db
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
