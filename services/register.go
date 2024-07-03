package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
	"github.com/jamesdavidyu/neighborhost-service/utils"
	"golang.org/x/crypto/bcrypt"
)

func CreateNeighbor(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var neighbor types.Register
		if err := json.NewDecoder(r.Body).Decode(&neighbor); err != nil {
			utils.WriteError(w, http.StatusBadRequest, err)
			return
		}

		if err := utils.Validate.Struct(neighbor); err != nil {
			errors := err.(validator.ValidationErrors)
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
			return
		}

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
