package service

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jamesdavidyu/neighborhost-service/types"
)

func GetTests(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM test")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		tests := []types.Test{}
		for rows.Next() {
			var t types.Test
			if err := rows.Scan(&t.ID, &t.Test); err != nil {
				log.Fatal(err)
			}
			tests = append(tests, t)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(tests)
	}
}
