package service

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jamesdavidyu/neighborhost-service/types"
)

func GetNeighborhoods(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM neighborhoods")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		neighborhoods := []types.Neighborhoods{}
		for rows.Next() {
			var n types.Neighborhoods
			if err := rows.Scan(&n.ID, &n.Neighborhood, &n.CreatedAt); err != nil {
				log.Fatal(err)
			}
			neighborhoods = append(neighborhoods, n)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(neighborhoods)
	}
}
