package controllers

import (
	"net/http"
	"strconv"

	"github.com/fpay/gopress"

	"blog/functions"
	"blog/models"
	"blog/services"
)

// PostController post controller
type PostController struct {
	db     *services.DBService
	title  string
	user   *models.User
	scRule *services.ScoreService
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
	c.scRule = app.Services.Get(services.ScoreServiceName).(*services.ScoreService)
	c.title = "BLOG-Article"
	app.GET("/posts", c.ListPosts)
	app.GET("/posts/create", c.CreatePage)
	app.GET("/posts/:id", c.ViewPost)
	app.POST("/posts/create", c.CreatePost)

}

// ListPosts Action
// Parameter gopress.Context is just alias of echo.Context
func (c *PostController) ListPosts(ctx gopress.Context) error {
	page := ctx.FormValue("page")
	orderBy := ctx.FormValue("sort")
	pageIndex, err := strconv.Atoi(page)
	if err != nil || page == "" {
		pageIndex = 1
	}

	if orderBy == "" {
		orderBy = "created_at desc"
	}

	limit := 10
	p := &models.Post{}
	pl, err := p.ListPosts(c.db.ORM, pageIndex, limit, orderBy)
	if err != nil {
		return ctx.Redirect(http.StatusFound, "/assets/404.html")
	}

	data := map[string]interface{}{
		"headTitle":    c.title,
		"avatar":       functions.GetAvatarURL(c.user.Avatar),
		"posts":        pl.Posts,
		"pagerContent": functions.GeneratePager(pl.Page, pl.Total, pl.Limit, pl.OrderBy, "/posts", nil),
	}

	return ctx.Render(http.StatusOK, "posts/list", data)
}

// CreatePage show create page
func (c *PostController) CreatePage(ctx gopress.Context) error {
	data := map[string]interface{}{
		"headTitle": c.title,
		"avatar":    functions.GetAvatarURL(c.user.Avatar),
	}
	cookie, err := ctx.Cookie("message")
	if err == nil {
		data["message"] = cookie.Value
		functions.SetCookieExpired(cookie)
		ctx.SetCookie(cookie)
	}

	return ctx.Render(http.StatusOK, "posts/create", data)
}

// CreatePost CreatePost
func (c *PostController) CreatePost(ctx gopress.Context) error {
	title := ctx.FormValue("title")
	content := ctx.FormValue("content")
	if title == "" || content == "" {
		ctx.SetCookie(functions.GetFlashCookie("message", "标题和内容不能为空"))
		return ctx.Redirect(http.StatusFound, "/posts/create")
	}
	if len(title) > 40 {
		ctx.SetCookie(functions.GetFlashCookie("message", "标题不能大于40个字符"))
		return ctx.Redirect(http.StatusFound, "/posts/create")
	}

	post := &models.Post{}
	err := post.CreatePost(c.db.ORM, c.user.ID, title, content, c.scRule.Rule.Post)
	if err != nil {
		ctx.SetCookie(functions.GetFlashCookie("message", err.Error()))
		return ctx.Redirect(http.StatusFound, "/posts/create")
	}

	return ctx.Redirect(http.StatusFound, "/posts")
}

// ViewPost view a post detail
func (c *PostController) ViewPost(ctx gopress.Context) error {
	idStr := ctx.Param("id")
	var postID uint64
	postID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return ctx.Redirect(http.StatusFound, "/assets/404.html")
	}

	post := &models.Post{}
	if c.db.ORM.First(post, postID).RecordNotFound() {
		return ctx.Redirect(http.StatusFound, "/assets/404.html")
	}

	comments := []*models.Comment{}
	comment := &models.Comment{}
	c.db.ORM.Model(post).Related(comment, "post_id").Find(&comments)

	data := map[string]interface{}{
		"headTitle": c.title,
		"post":      post,
		"comments":  comments,
		"avatar":    functions.GetAvatarURL(c.user.Avatar),
	}

	return ctx.Render(http.StatusOK, "posts/detail", data)
}
