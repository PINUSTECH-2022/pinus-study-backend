package sql

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func SqlDirExec(db *sql.DB, dir string) {
	fileList, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range fileList {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			filePath := filepath.Join(dir, file.Name())
			query, err := ioutil.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file %s: %s\n", filePath, err)
				continue
			}

			_, err = db.Exec(string(query))
			if err != nil {
				fmt.Printf("Error executing file %s: %s\n", filePath, err)
				continue
			}

			fmt.Printf("SQL file %s executed successfully\n", filePath)
		}
	}
}

func ExecSQLFiles(db *sql.DB, sqlDir string, sqlFiles []string) {
	for _, sqlFile := range sqlFiles {
		query, err := ioutil.ReadFile(sqlDir + sqlFile)
		if err != nil {
			panic(err)
		}

		_, err1 := db.Exec(string(query))
		if err1 != nil {
			panic(err1)
		}

		fmt.Printf("%s%s has been run successfully\n", sqlDir, sqlFile)
	}
}
