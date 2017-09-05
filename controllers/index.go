package controllers

import (
	"net/http"

	"github.com/fpay/gopress"
	"github.com/jinzhu/gorm"
	"blog/services"
	"blog/models"
)

// IndexController
type IndexController struct {
	// Uncomment this line if you want to use services in the app
	// app *gopress.App
	db *services.DBService
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
	app.GET("/", c.HomeAction)
	// app.POST("/index/sample", c.SamplePostAction)
	// app.PUT("/index/sample", c.SamplePutAction)
	// app.DELETE("/index/sample", c.SampleDeleteAction)
}

// SampleGetAction Action
// Parameter gopress.Context is just alias of echo.Context
func (c *IndexController) HomeAction(ctx gopress.Context) error {
	// Or you can get app from request context
	// app := gopress.AppFromContext(ctx)

	return ctx.String(http.StatusOK, "welcome to my blog")
}
