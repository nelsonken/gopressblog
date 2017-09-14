package main

import (
	"blog/config"
	"blog/controllers"
	"blog/middlewares"
	"blog/services"

	"github.com/fpay/gopress"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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
	opts := &config.BlogOptions{}
	opts.Database = &services.DBOptions{}
	opts.ScoreRule = &services.ScoreRule{}
	opts.Elastic = &services.ElasticOption{}
	config.GetConfig(ConfigFile, opts)

	// services register
	dbs := services.NewDBService(opts.Database.DBType, opts.Database)
	vs := services.NewValidatorService()
	score := services.NewScoreService(opts.ScoreRule)
	es := services.NewElasticService(opts.Elastic)
	s.RegisterServices(dbs, vs, score, es)

	// register middlewares
	s.RegisterGlobalMiddlewares(
		gopress.NewLoggingMiddleware("global", gopress.NewLogger()),
	)

	// RouteGroups route groups
	needLoginMiddlewares := []gopress.MiddlewareFunc{middlewares.NewAuthMiddleware()}

	authGroup := s.App().Group("/blog", needLoginMiddlewares...)
	//init and register controllers
	s.RegisterControllers(
		controllers.NewIndexController(),
		controllers.NewUserController(),
		controllers.NewPostController(authGroup),
		controllers.NewCommentController(authGroup),
		controllers.NewAccountController(authGroup),
	)

	// static path
	s.App().Static("/assets", "assets")
	//
	s.Start()
}
