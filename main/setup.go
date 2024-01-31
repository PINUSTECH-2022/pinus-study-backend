package main

import (
	"example/web-service-gin/database"
	"example/web-service-gin/imports"
	"example/web-service-gin/sql"
	"fmt"
)

func main() {
	db := database.GetDb()

	sql.SetupTable(db)
	sql.SetupTrigger(db)
	sql.SetupProcedure(db)

	fmt.Println("Database has been setup successfully")

	imports.ModulesSetup()
}
