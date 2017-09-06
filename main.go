package main

import (
	"blog/config"
	"blog/controllers"
	"blog/middlewares"
	"blog/models"
	"blog/services"

	"github.com/fpay/gopress"
)

const (
	// ConfigFile config file path
	ConfigFile = "config/config.yaml"
	// TimeFormat time format str
	TimeFormat = "2006-01-02 15:04:05"
)

func main() {
	// create server
	s := gopress.NewServer(gopress.ServerOptions{
		Port: 3000,
	})

	opts := &config.Options{}
	opts.Database = &services.DBOptions{}
	opts.ScoreRule = &services.ScoreRule{}
	config.GetConfig(ConfigFile, opts)

	dbs := services.NewDBService(opts.Database.DBType, opts.Database)
	vs := services.NewValidatorService()
	us := services.NewUserService()
	us.User = &models.User{}
	score := services.NewScoreService(opts.ScoreRule)
	s.RegisterServices(dbs, vs, us, score)
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
