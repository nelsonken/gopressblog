package controllers

import (
	"net/http"

	"github.com/fpay/gopress"
	"blog/services"
	"blog/models"
	. "blog/functions"
)

// IndexController
type IndexController struct {
	// Uncomment this line if you want to use services in the app
	app *gopress.App
	db *services.DBService
	currentUser *models.User
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
	c.title = "Home"
	c.currentUser = app.Services.Get(services.UserServiceName).(*services.UserService).User
	app.GET("/", c.Home)
}

// HomeAction Action
// show some no use data analyes
// Parameter gopress.Context is just alias of echo.Context
func (c *IndexController) Home(ctx gopress.Context) error {
	data := map[string]interface{}{
		"titile": c.title,
		"avatar": GetAvatarURL(c.currentUser.Avatar),
	}

	return ctx.Render(http.StatusOK, "index", data)
}
