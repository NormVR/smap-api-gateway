package api

import (
	"api-gateway/internal/handler"
	"api-gateway/internal/services"
	"context"
	"log"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HttpServer struct {
	Server *http.Server
}

func NewHttpServer() *HttpServer {
	userService := services.NewUserService()
	userHandler := handler.NewUserHandler(userService)

	publicRouter := http.NewServeMux()
	protectedRouter := http.NewServeMux()

	publicRouter.HandleFunc("POST /auth/register", userHandler.Register)
	publicRouter.HandleFunc("POST /auth/login", userHandler.Login)
	publicRouter.HandleFunc("POST /auth/logout", userHandler.Logout)

	mainRouter := http.NewServeMux()
	mainRouter.Handle("/auth/", publicRouter)
	mainRouter.Handle("/api/", authMiddleware(protectedRouter))

	return &HttpServer{
		Server: &http.Server{
			Addr:    ":8080",
			Handler: mainRouter,
		},
	}
}

func (s *HttpServer) RunServer() error {
	log.Printf("Starting http server on %s\n", s.Server.Addr)
	if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	log.Printf("Http server started\n")
	return nil
}

func (s *HttpServer) StopServer(ctx context.Context) {
	log.Printf("Stopping http server\n")
	if err := s.Server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v\n", err)
		if err = s.Server.Close(); err != nil {
			log.Printf("Server force shutdown error: %v\n", err)
		}
		return
	}
	log.Printf("Http server stopped\n")
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		if tokenString == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		userService := services.NewUserService()
		userId, err := userService.ValidateToken(tokenString)

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

		if userId == 0 {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
