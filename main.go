package main

import (
	"blog/config"
	"blog/controllers"
	"blog/services"
	"github.com/fpay/gopress"
)

const (
	ConfigFile = "config/config.yaml"
)

func main() {
	// create server
	s := gopress.NewServer(gopress.ServerOptions{
		Port: 3000,
	})

	opts := &config.Options{}
	opts.Database = &services.DBOptions{}
	config.GetConfig(ConfigFile, opts)

	dbService := services.NewDBService(opts.Database.DBType, opts.Database)
	s.RegisterServices(dbService)

	// register middlewares
	s.RegisterGlobalMiddlewares(
		gopress.NewLoggingMiddleware("global", gopress.NewLogger()),
	)

	//init and register controllers
	s.RegisterControllers(
		controllers.NewIndexController(),
	)

	//
	s.Start()
}
