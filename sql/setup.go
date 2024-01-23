package main

import (
	"example/web-service-gin/database"
	"example/web-service-gin/imports"
	"fmt"
	"io/ioutil"
)

func main() {
	db := database.GetDb()

	sqlFile := "sql/setup.sql"

	query, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		panic(err)
	}

	_, err1 := db.Exec(string(query))
	if err1 != nil {
		panic(err1)
	}

	fmt.Println("Database has been setup successfully")

	imports.ModulesSetup()
}
