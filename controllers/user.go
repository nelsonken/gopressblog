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
	app   *gopress.App
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
	c.app = app
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
		ctx.SetCookie(functions.GetFlashCookie("message", "账号与密码不能为空"))
		return ctx.Redirect(http.StatusFound, "/login")
	}

	u := &models.User{}
	err := u.Login(c.db.ORM, name, password)
	if err != nil {
		ctx.SetCookie(functions.GetFlashCookie("message", err.Error()))
		return ctx.Redirect(http.StatusFound, "/login")
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
	}
	cookie, err := ctx.Cookie("message")
	if err == nil {
		data["message"] = cookie.Value
		ctx.SetCookie(functions.SetCookieExpired(cookie))
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
		ctx.SetCookie(functions.GetFlashCookie("message", err.Error()))
		return ctx.Redirect(http.StatusFound, "/register")
	}

	if rf.Agree != "agree" {
		ctx.SetCookie(functions.GetFlashCookie("message", "请阅读注册协议"))
		return ctx.Redirect(http.StatusFound, "/register")
	}

	u := &models.User{}
	err := u.Register(c.db.ORM, rf.Name, rf.Password)

	if err != nil {
		ctx.SetCookie(functions.GetFlashCookie("message", err.Error()))
		return ctx.Redirect(http.StatusFound, "/register")
	}

	return ctx.Redirect(http.StatusFound, "/login")
}

// RegisterPage Action
// Parameter gopress.Context is just alias of echo.Context
func (c *UserController) RegisterPage(ctx gopress.Context) error {
	data := map[string]interface{}{
		"headTitle": "注册",
	}
	cookie, err := ctx.Cookie("message")
	if err == nil {
		data["message"] = cookie.Value
		ctx.SetCookie(functions.SetCookieExpired(cookie))
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
