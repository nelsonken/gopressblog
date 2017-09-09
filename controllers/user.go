package controllers

import (
	"net/http"
	"strconv"
	"time"

	"blog/functions"
	"blog/models"
	"blog/requests"
	"blog/services"

	"github.com/fpay/gopress"
)

// UserController user controller
type UserController struct {
	db    *services.DBService
	valid *services.Validator
}

// NewUserController returns index controller instance.
func NewUserController() *UserController {
	return new(UserController)
}

// RegisterRoutes registes routes to app
// It is used to implements gopress.Controller.
func (c *UserController) RegisterRoutes(app *gopress.App) {
	c.db = app.Services.Get(services.DBServerName).(*services.DBService)
	c.valid = app.Services.Get(services.ValidatorServiceName).(*services.Validator)
	app.GET("/logout", c.Logout)
	app.GET("/login", c.LoginPage)
	app.POST("/login", c.Login)
	app.GET("/register", c.RegisterPage)
	app.POST("/register", c.Register)
}

// Login Action
// Parameter gopress.Context is just alias of echo.Context
func (c *UserController) Login(ctx gopress.Context) error {
	name := ctx.FormValue("name")
	password := ctx.FormValue("password")
	autoLogin := ctx.FormValue("autoLogin")
	if name == "" || password == "" {
		return ctx.Redirect(http.StatusFound, "/login?message=账号密码不能为空")
	}

	u := &models.User{}
	err := u.Login(c.db.ORM, name, password)
	if err != nil {
		return ctx.Redirect(http.StatusFound, "/login?message="+err.Error())
	}
	var expired time.Time
	if autoLogin == "yes" {
		expired = time.Now().Add(time.Minute * 30)
	} else {
		expired = time.Now().Add(time.Hour * 72)
	}
	cookie := &http.Cookie{Name: "uid", Value: strconv.FormatUint(uint64(u.ID), 10), Expires: expired}
	ctx.SetCookie(cookie)

	return ctx.Redirect(http.StatusFound, "/")
}

// LoginPage Action
// Parameter gopress.Context is just alias of echo.Context
func (c *UserController) LoginPage(ctx gopress.Context) error {
	data := map[string]interface{}{
		"headTitle": "登录",
		"message":   ctx.QueryParam("message"),
	}

	return ctx.Render(http.StatusOK, "user/login", data)
}

// Register Action
// Parameter gopress.Context is just alias of echo.Context
func (c *UserController) Register(ctx gopress.Context) error {
	rf := &requests.RegisterForm{
		Name:            ctx.FormValue("name"),
		Password:        ctx.FormValue("password"),
		PasswordConfirm: ctx.FormValue("password_confirm"),
		Agree:           ctx.FormValue("agree"),
	}
	if err := c.valid.Validate(rf); err != nil {
		return ctx.Redirect(http.StatusFound, "/register?message="+err.Error())
	}

	if rf.Agree != "agree" {
		return ctx.Redirect(http.StatusFound, "/register?message=请阅读注册协议")
	}

	u := &models.User{}
	err := u.Register(c.db.ORM, rf.Name, rf.Password)

	if err != nil {
		return ctx.Redirect(http.StatusFound, "/register?message="+err.Error())
	}

	return ctx.Redirect(http.StatusFound, "/login")
}

// RegisterPage Action
// Parameter gopress.Context is just alias of echo.Context
func (c *UserController) RegisterPage(ctx gopress.Context) error {
	data := map[string]interface{}{
		"headTitle": "注册",
		"message":   ctx.QueryParam("message"),
	}

	return ctx.Render(http.StatusOK, "user/register", data)
}

// Logout Action
// Parameter gopress.Context is just alias of echo.Context
func (c *UserController) Logout(ctx gopress.Context) error {
	cookie, err := ctx.Cookie("uid")
	if err == nil {
		cookie = &http.Cookie{Name: "uid", Value: ""}
	}

	ctx.SetCookie(functions.SetCookieExpired(cookie))

	return ctx.Redirect(http.StatusFound, "/login")
}
