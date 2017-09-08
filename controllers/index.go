package controllers

import (
	"net/http"

	"blog/services"

	"github.com/fpay/gopress"
)

// IndexController action pointer
type IndexController struct {
	// Uncomment this line if you want to use services in the app
	app   *gopress.App
	db    *services.DBService
	title string
}

// NewIndexController returns index controller instance.
func NewIndexController() *IndexController {
	return new(IndexController)
}

// RegisterRoutes registes routes to app
// It is used to implements gopress.Controller.
func (c *IndexController) RegisterRoutes(app *gopress.App) {
	// Uncomment this line if you want to use services in the app
	// c.app = app
	c.db = app.Services.Get(services.DBServerName).(*services.DBService)
	c.app = app
	c.title = "首页"
	app.GET("/", c.Home)
}

// Home Action
// show some no use data analyes
// Parameter gopress.Context is just alias of echo.Context
func (c *IndexController) Home(ctx gopress.Context) error {
	return ctx.Redirect(http.StatusMovedPermanently, "/blog/posts")
}
