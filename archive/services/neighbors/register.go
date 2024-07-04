package neighbors

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	"github.com/go-playground/validator/v10"
// 	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
// 	"github.com/jamesdavidyu/neighborhost-service/utils"
// 	"golang.org/x/crypto/bcrypt"
// )

// type Handler struct {
// 	store types.NeighborStore
// }

// func NewHandler(store types.NeighborStore) *Handler {
// 	return &Handler{store: store}
// }

// func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
// 	var neighbor types.Register
// 	if err := json.NewDecoder(r.Body).Decode(&neighbor); err != nil {
// 		utils.WriteError(w, http.StatusBadRequest, err)
// 		return
// 	}

// 	if err := utils.Validate.Struct(neighbor); err != nil {
// 		errors := err.(validator.ValidationErrors)
// 		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
// 		return
// 	}

// 	// need to check if user already exists

// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(neighbor.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	err = h.store.CreateNeighbor(types.Neighbors{
// 		Email:    neighbor.Email,
// 		Username: neighbor.Username,
// 		Zipcode:  neighbor.Zipcode,
// 		Password: string(hashedPassword),
// 	})

// 	if err != nil {
// 		utils.WriteError(w, http.StatusInternalServerError, err)
// 		return
// 	}

// 	w.WriteHeader(http.StatusCreated)
// 	json.NewEncoder(w).Encode(neighbor)
// }
