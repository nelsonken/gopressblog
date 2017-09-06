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
)

func main() {
	// create server
	s := gopress.NewServer(gopress.ServerOptions{
		Port: 3000,
	})

	// opt
	opts := &config.Options{}
	opts.Database = &services.DBOptions{}
	opts.ScoreRule = &services.ScoreRule{}
	config.GetConfig(ConfigFile, opts)

	// services register
	dbs := services.NewDBService(opts.Database.DBType, opts.Database)
	vs := services.NewValidatorService()
	us := services.NewUserService()
	us.User = &models.User{}
	score := services.NewScoreService(opts.ScoreRule)
	s.RegisterServices(dbs, vs, us, score)

	// register middlewares
	s.RegisterGlobalMiddlewares(
		gopress.NewLoggingMiddleware("global", gopress.NewLogger()),
	)

	//init and register controllers
	s.RegisterControllers(
		controllers.NewIndexController(),
		controllers.NewUserController(),
		controllers.NewPostController(),
		controllers.NewCommentController(),
	)

	s.App().Group("/blog").Use(middlewares.NewAuthMiddleware(us.User))

	// static path
	s.App().Static("/assets", "assets")
	//
	s.Start()
}
