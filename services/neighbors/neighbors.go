package neighbors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/jamesdavidyu/neighborhost-service/cmd/model/types"
	"github.com/jamesdavidyu/neighborhost-service/config"
	"github.com/jamesdavidyu/neighborhost-service/services/auth"
	"github.com/jamesdavidyu/neighborhost-service/utils"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store types.NeighborStore
}

func NewHandler(store types.NeighborStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/auth/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/auth/login", h.handleLogin).Methods("POST")
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
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

	// need to check if user already exists

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(neighbor.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateNeighbor(types.Neighbors{
		Email:    neighbor.Email,
		Username: neighbor.Username,
		Zipcode:  neighbor.Zipcode,
		Password: string(hashedPassword),
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(neighbor)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var login types.Login
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(login); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
		return
	}

	neighbor, err := h.store.GetNeighborWithEmailOrUsername(login.EmailOrUsername)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found"))
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(neighbor.Password), []byte(login.Password)) != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found"))
	} else {
		secret := []byte(config.Envs.JWTSecret)
		token, err := auth.CreateJWT(secret, neighbor.Id)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		json.NewEncoder(w).Encode(map[string]any{
			"token":          token,
			"neighborId":     neighbor.Id,
			"email":          neighbor.Email,
			"username":       neighbor.Username,
			"zipcode":        neighbor.Zipcode,
			"neighborhoodId": neighbor.NeighborhoodID,
			"verified":       neighbor.Verified,
		})
	}
}
