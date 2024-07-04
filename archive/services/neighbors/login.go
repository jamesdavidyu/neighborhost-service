package neighbors

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/go-playground/validator/v10"
// 	"github.com/jamesdavidyu/neighborhost-service/cmd/model/db"
// 	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
// 	"github.com/jamesdavidyu/neighborhost-service/config"
// 	"github.com/jamesdavidyu/neighborhost-service/services/auth"
// 	"github.com/jamesdavidyu/neighborhost-service/utils"
// 	"golang.org/x/crypto/bcrypt"
// )

// func Login(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var login types.Login
// 		if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
// 			utils.WriteError(w, http.StatusBadRequest, err)
// 			return
// 		}

// 		if err := utils.Validate.Struct(login); err != nil {
// 			errors := err.(validator.ValidationErrors)
// 			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
// 			return
// 		}

// 		neighbor, err := GetNeighborWithEmailOrUsername(login.EmailOrUsername)
// 		if err != nil {
// 			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found"))
// 			return
// 		}

// 		if bcrypt.CompareHashAndPassword([]byte(neighbor.Password), []byte(login.Password)) != nil {
// 			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found"))
// 		} else {
// 			secret := []byte(config.Envs.JWTSecret)
// 			token, err := auth.CreateJWT(secret, neighbor.ID)
// 			if err != nil {
// 				utils.WriteError(w, http.StatusInternalServerError, err)
// 				return
// 			}

// 			json.NewEncoder(w).Encode(map[string]string{"token": token})
// 		}
// 	}
// }
