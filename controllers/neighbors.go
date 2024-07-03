package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
	"github.com/jamesdavidyu/neighborhost-service/utils"
	"golang.org/x/crypto/bcrypt"
)

func CreateNeighbor(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var neighbor types.Register
		json.NewDecoder(r.Body).Decode(&neighbor)

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(neighbor.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		_, err = db.Exec(`INSERT INTO neighbors (email, username, zipcode, password)
							VALUES ($1, $2, $3, $4)`, neighbor.Email, neighbor.Username, neighbor.Zipcode, hashedPassword)
		if err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(neighbor)
	}
}
