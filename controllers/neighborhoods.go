package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
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
			var neighborhood types.Neighborhoods
			if err := rows.Scan(&neighborhood.ID, &neighborhood.Neighborhood, &neighborhood.CreatedAt); err != nil {
				log.Fatal(err)
			}
			neighborhoods = append(neighborhoods, neighborhood)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(neighborhoods)
	}
}

func CreateNeighborhood(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var neighborhood types.Neighborhoods
		json.NewDecoder(r.Body).Decode(&neighborhood)

		_, err := db.Exec(`INSERT INTO neighborhoods (neighborhood)
							VALUES ($1)`, neighborhood.Neighborhood)
		if err != nil {
			fmt.Println(err)
		}

		json.NewEncoder(w).Encode(neighborhood)
	}
}
