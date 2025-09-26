package handler

import (
	"api-gateway/internal/models/auth"
	"api-gateway/internal/services"
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	err := json.NewDecoder(r.Body).Decode(&userData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	/*_, err = mail.ParseAddress(userData.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}*/

	err = h.service.RegisterUser(&userData)

	if err != nil {
		log.Printf("failed to register user: %v", err)
		st, ok := status.FromError(err)
		if !ok {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		switch st.Code() {
		case codes.InvalidArgument:
			http.Error(w, st.Message(), http.StatusBadRequest)
			return
		case codes.AlreadyExists:
			http.Error(w, st.Message(), http.StatusConflict)
			return
		case codes.Internal:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
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

	defer r.Body.Close()

	result, err := h.service.LoginUser(&loginData)

	if err != nil {
		log.Printf("failed to login user: %v", err)
		st, ok := status.FromError(err)
		if !ok {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		switch st.Code() {
		case codes.InvalidArgument:
			http.Error(w, st.Message(), http.StatusBadRequest)
			return
		case codes.Unauthenticated:
			http.Error(w, st.Message(), http.StatusUnauthorized)
			return
		case codes.Internal:
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(result)
}

func (h *UserHandler) Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}
