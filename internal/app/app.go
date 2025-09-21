package app

import (
	"api-gateway/internal/api"
	"api-gateway/internal/configs"
	"context"
)

type App struct {
	IntegrationConfig *configs.IntegrationConfig
	ServiceConfig     *configs.ServicesConfig
	HttpServer        *api.HttpServer
}

func New() *App {
	return &App{
		IntegrationConfig: configs.NewIntegrationConfig(),
		ServiceConfig:     configs.NewServiceConfig(),
		HttpServer:        api.NewHttpServer(),
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
