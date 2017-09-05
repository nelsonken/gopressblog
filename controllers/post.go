package controllers

import (
	"strconv"

	"github.com/fpay/gopress"
	"blog/services"
	"blog/models"
	. "blog/functions"
	"net/http"
)

// PostController
type PostController struct {
	db    *services.DBService
	title string
	user *models.User
}

// NewPostController returns post controller instance.
func NewPostController() *PostController {
	return new(PostController)
}

// RegisterRoutes registes routes to app
// It is used to implements gopress.Controller.
func (c *PostController) RegisterRoutes(app *gopress.App) {
	c.db = app.Services.Get(services.DBServerName).(*services.DBService)
	c.user = app.Services.Get(services.UserServiceName).(*services.UserService).User
	c.title = "文章列表"
	app.GET("/posts", c.ListPosts)

}

// ListPosts Action
// Parameter gopress.Context is just alias of echo.Context
func (c *PostController) ListPosts(ctx gopress.Context) error {
	page := ctx.FormValue("page")
	orderBy := ctx.FormValue("sort")
	pageIndex, err := strconv.Atoi(page)
	if err != nil {
		pageIndex = 1
	}

	if orderBy == "" {
		orderBy = "create_at desc"
	}

	limit := 10
	p := &models.Post{}
	pl, err := p.ListPosts(c.db.ORM, pageIndex, limit, orderBy)
	if err != nil {
		return ctx.String(200, err.Error())
		//return ctx.Redirect(301, "/assets/404.html")
	}

	data := map[string]interface{}{
		"title":        c.title,
		"posts":        pl.Posts,
		"pagerContent": GeneratePager(pl.Page, pl.Total, pl.Limit, pl.OrderBy, "/posts", nil),
	}

	return ctx.Render(http.StatusOK, "posts", data)
}
