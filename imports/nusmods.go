package imports

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type modules struct {
	Title       string
	ModuleCode  string
	Description string
}

func GetDbNusMods() *sql.DB {
	dbUri := os.Getenv("DATABASE_URI")

	if dbUri == "" {
		err := godotenv.Load(".env")

		if err != nil {
			panic(err)
		}

		dbUri = os.Getenv("DATABASE_URI")
	}

	db, err := sql.Open("postgres", dbUri)
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

func ModulesSetup() {
	API_BASE_URL := "https://api.nusmods.com/v2"
	acadYear := "2023-2024" // must be in yyyy-yyyy format
	fetchUrl := fmt.Sprintf("%s/%s/%s", API_BASE_URL, acadYear, "moduleInfo.json")
	resp, err := http.Get(fetchUrl)

	if err != nil {
		panic(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	resp.Body.Close()
	var mods []modules
	jsonErr := json.Unmarshal(body, &mods)

	if jsonErr != nil {
		panic(jsonErr)
	}

	db := GetDbNusMods()
	defer db.Close()

	tx, dbErr := db.Begin()

	if dbErr != nil {
		panic(dbErr)
	}

	stmt, prepErr := tx.Prepare("INSERT INTO Modules VALUES($1, $2, $3) ON CONFLICT DO NOTHING;")

	if prepErr != nil {
		panic(prepErr)
	}

	for i := 0; i < len(mods); i++ {
		fmt.Printf("%d %s\n", i, mods[i].ModuleCode)
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

func main() {
	ModulesSetup()
}
