package main

import (
	"blog/config"
	"blog/controllers"
	"blog/services"
	"github.com/fpay/gopress"
	"blog/middlewares"
	"blog/models"
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

	dbs := services.NewDBService(opts.Database.DBType, opts.Database)
	vs := services.NewVlidatorService()
	us := services.NewUserService()
	us.User = &models.User{}
	s.RegisterServices(dbs, vs, us)
	// register middlewares
	s.RegisterGlobalMiddlewares(
		gopress.NewLoggingMiddleware("global", gopress.NewLogger()),
		middlewares.NewAuthMiddleware(us.User),
	)

	//init and register controllers
	s.RegisterControllers(
		controllers.NewIndexController(),
		controllers.NewUserController(),
		controllers.NewPostController(),
	)
	s.App().Static("/assets", "assets")
	//
	s.Start()
}
