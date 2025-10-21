package api

import (
	"api-gateway/internal/handler"
	"context"
	"log"
	"net/http"
)

type HttpServer struct {
	Server      *http.Server
	UserHandler *handler.UserHandler
}

func NewHttpServer(userHandler *handler.UserHandler) *HttpServer {

	publicRouter := http.NewServeMux()
	protectedRouter := http.NewServeMux()

	publicRouter.HandleFunc("POST /auth/register", userHandler.Register)
	publicRouter.HandleFunc("POST /auth/login", userHandler.Login)
	publicRouter.HandleFunc("POST /auth/logout", userHandler.Logout)

	protectedRouter.HandleFunc("GET /api/users/{id}", userHandler.GetUser)
	//protectedRouter.HandleFunc("PUT /api/users/{id}", userHandler.UpdateUser)
	//protectedRouter.HandleFunc("POST /api/users/{id}/follow{follower_id}", userHandler.Follow)
	//protectedRouter.HandleFunc("GET /api/users/search", userHandler.Search)

	mainRouter := http.NewServeMux()
	mainRouter.Handle("/auth/", publicRouter)
	mainRouter.Handle("/api/", userHandler.AuthMiddleware(protectedRouter))

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
