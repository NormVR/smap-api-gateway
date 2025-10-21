package handler

import (
	"api-gateway/internal/models/auth"
	"api-gateway/internal/services"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
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

	userData := auth.AuthData{}

	err := json.NewDecoder(r.Body).Decode(&userData)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

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
		default:
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusInternalServerError)
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

	loginData := auth.AuthData{}

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
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(result)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	tokenString := r.Header.Get("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	if tokenString == "" {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	err := h.service.Logout(tokenString)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	path := r.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 4 {
		http.Error(w, "incorrect URL", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(parts[3])

	if err != nil {
		http.Error(w, "incorrect URL", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(id)

	if err != nil {
		log.Printf("failed to get user: %v", err)
		st, ok := status.FromError(err)
		if !ok {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		switch st.Code() {
		case codes.NotFound:
			http.Error(w, st.Message(), http.StatusNotFound)
			return
		case codes.Internal:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		if tokenString == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		userId, err := h.service.ValidateToken(tokenString)

		if err != nil {
			log.Printf("Error validating token: %v\n", err)

			st, ok := status.FromError(err)

			if !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusInternalServerError)
				return
			}

			switch st.Code() {
			case codes.InvalidArgument:
				http.Error(w, st.Message(), http.StatusBadRequest)
				return
			case codes.Unauthenticated:
				http.Error(w, st.Message(), http.StatusUnauthorized)
				return
			default:
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusInternalServerError)
				return
			}
		}

		if userId == uuid.Nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
