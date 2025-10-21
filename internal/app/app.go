package app

import (
	"api-gateway/internal/api"
	"api-gateway/internal/configs"
	"api-gateway/internal/grpc/auth"
	"api-gateway/internal/handler"
	"api-gateway/internal/repository/redis"
	"api-gateway/internal/services"
	"context"
)

type App struct {
	HttpServer *api.HttpServer
}

func New() *App {

	grpcClient := auth.New(configs.NewServiceConfig())
	redisRepo := redis.NewRedisRepository(configs.NewIntegrationConfig())
	userService := services.NewUserService(grpcClient, redisRepo)
	userHandler := handler.NewUserHandler(userService)
	httpServer := api.NewHttpServer(userHandler)

	return &App{
		HttpServer: httpServer,
	}
}

func (app *App) Run(serverErr chan error) {
	if err := app.HttpServer.RunServer(); err != nil {
		serverErr <- err
	}
	close(serverErr)
}

func (app *App) Stop(ctx context.Context) {
	app.HttpServer.StopServer(ctx)
}
