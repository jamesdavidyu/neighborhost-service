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
	var register types.Register
	if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(register); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid submission for %v", errors))
		return
	}

	// need to handle errors better
	checkEmail, err := h.store.GetNeighborWithEmail(register.Email)
	if checkEmail == nil {
		checkUsername, err := h.store.GetNeighborWithUsername(register.Username)
		if checkUsername == nil {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			err = h.store.CreateNeighbor(types.Neighbors{
				Email:    register.Email,
				Username: register.Username,
				Zipcode:  register.Zipcode,
				Password: string(hashedPassword),
			})

			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			} else {
				neighbor, err := h.store.GetNeighborWithEmailOrUsername(register.Email)
				if err != nil {
					utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found"))
					return
				}

				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(neighbor) // need to return token and ID? Need to run getNeighborById again?
			}
		} else {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("taken"))
		}
		if err != nil {
			return
		}
	} else {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("taken"))
	}
	if err != nil {
		return
	}
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
			"neighborhoodId": neighbor.NeighborhoodId,
			"verified":       neighbor.Verified,
		})
	}
}
