package users

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	user, err := h.service.createUser(r.Context(), req)
	if err != nil {
		if errors.Is(err, errUsernameTaken) {
			http.Error(w, errUsernameTaken.Error(), http.StatusConflict)
			return

		}

		log.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// func GetUser(w http.ResponseWriter, r *http.Request) {
// 	id, err := strconv.Atoi(r.PathValue("id"))

// 	if err != nil {
// 		http.Error(w, "Invalid id", http.StatusBadRequest)
// 		return
// 	}

// 	user := repository.getUserById(id)

// 	if user == nil {
// 		http.Error(w, "Not found", http.StatusNotFound)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(user)
// }
