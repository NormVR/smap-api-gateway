package handler

import (
	"api-gateway/internal/models/auth"
	"api-gateway/internal/services"
	"encoding/json"
	"net/http"
	"net/mail"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userData := auth.UserData{}

	err := json.NewDecoder(r.Body).Decode(userData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	_, err = mail.ParseAddress(userData.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = h.service.RegisterUser(&userData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	loginData := auth.LoginData{}

	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	result, err := h.service.LoginUser(&loginData)

	json.NewEncoder(w).Encode(result)
}
