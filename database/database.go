package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func db() {
	DATABASE_URL := "postgres://ddbrflhccgsbnv:9dea2c78cfd420c6af3f51e1bac56c9a15ed73dc73ad1b50172c5e973b613e37@ec2-44-206-137-96.compute-1.amazonaws.com:5432/devchlvlplpfpu"
	db, err := sql.Open("postgres", os.Getenv(DATABASE_URL))
	if err != nil {
		panic(err)
	}
}
