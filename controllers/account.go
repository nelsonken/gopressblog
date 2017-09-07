package controllers

import (
	"blog/functions"
	"blog/models"
	"blog/services"
	"net/http"
	"strconv"

	"github.com/fpay/gopress"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

// AccountController actions is here
type AccountController struct {
	group *echo.Group
	orm   *gorm.DB
	user  *models.User
	title string
}

// NewAccountController returns account controller instance.
func NewAccountController(group *echo.Group) *AccountController {
	c := new(AccountController)
	c.group = group
	group.GET("/messages", c.ListMessages)
	c.title = "消息"
	return c
}

// RegisterRoutes registes routes to app
// It is used to implements gopress.Controller.
func (c *AccountController) RegisterRoutes(app *gopress.App) {
	c.orm = app.Services.Get(services.DBServerName).(*services.DBService).ORM
	c.user = app.Services.Get(services.UserServiceName).(*services.UserService).User
}

// ListMessages Action
// Parameter gopress.Context is just alias of echo.Context
func (c *AccountController) ListMessages(ctx gopress.Context) error {
	p := ctx.QueryParam("page")
	sortBy := ctx.QueryParam("sort")
	page, _ := strconv.Atoi(p)
	if page == 0 {
		page = 1
	}
	m := &models.Message{}
	msgs := []*models.Message{}
	limit := 10
	if sortBy == "" {
		sortBy = defaultSortBy
	}
	total := m.ListMessages(c.orm, c.user.ID, &msgs, limit, page, sortBy)

	data := map[string]interface{}{
		"avatar":    functions.GetAvatarURL(c.user.Avatar),
		"headTitle": c.title,
		"msgs":      msgs,
		"pager":     functions.GeneratePager(page, total, limit, sortBy, "/messages", nil),
	}

	return ctx.Render(http.StatusOK, "account/message", data)
}
