package main

import (
	"pinus/database"
	"fmt"
	"io"
	_ "github.com/lib/pq"
	"net/http"
	"encoding/json"
)

type modules struct {
	Title string 
	ModuleCode string
	Description string
}

func main() {
	API_BASE_URL := "https://api.nusmods.com/v2"
	acadYear := "2022-2023" // must be in yyyy-yyyy format
	fetchUrl := fmt.Sprintf("%s/%s/%s", API_BASE_URL, acadYear, "moduleInfo.json")
	resp, err := http.Get(fetchUrl)

	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	var mods []modules
	jsonErr := json.Unmarshal(body, &mods)

	if jsonErr != nil {
		panic(jsonErr)
	}

	db := database.GetDb()
	defer db.Close()

	tx, dbErr := db.Begin()

	if dbErr != nil {
		panic(dbErr)
	}

	stmt, prepErr := tx.Prepare("INSERT INTO Modules VALUES($1, $2, $3) ON CONFLICT DO NOTHING")

	if prepErr != nil {
		panic(prepErr)
	}

	for i := 0; i < len(mods); i++ {
		fmt.Printf("%d\n", i)
		_, stmtErr := stmt.Exec(mods[i].ModuleCode, mods[i].Title, mods[i].Description)

		if stmtErr != nil {
			panic(stmtErr)
		}
	}

	commitErr := tx.Commit()

	if commitErr != nil {
		panic(commitErr)
	}

	stmt.Close()
}
