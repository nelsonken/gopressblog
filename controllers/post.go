package controllers

import (
	"net/http"
	"strconv"

	"github.com/fpay/gopress"

	"blog/functions"
	"blog/models"
	"blog/services"
	"fmt"
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
	app.GET("/blog/posts", c.ListPosts)
	app.GET("/blog/posts/create", c.CreatePage)
	app.GET("/blog/posts/update/:id", c.UpdatePage)
	app.POST("/blog/posts/update", c.UpdatePost)
	app.GET("/blog/posts/:id", c.ViewPost)
	app.POST("/blog/posts/create", c.CreatePost)
	app.GET("/blog/myposts", c.MyPosts)

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
		"pagerContent": functions.GeneratePager(pl.Page, pl.Total, pl.Limit, pl.OrderBy, "/blog/posts", nil),
		"getAuthorName": func(id uint) string {
			return ""
		},
	}

	return ctx.Render(http.StatusOK, "posts/list", data)
}

// CreatePage show create page
func (c *PostController) CreatePage(ctx gopress.Context) error {
	data := map[string]interface{}{
		"headTitle": c.title,
		"avatar":    functions.GetAvatarURL(c.user.Avatar),
		"message":   ctx.QueryParam("message"),
	}

	return ctx.Render(http.StatusOK, "posts/create", data)
}

// CreatePost CreatePost
func (c *PostController) CreatePost(ctx gopress.Context) error {
	title := ctx.FormValue("title")
	content := ctx.FormValue("content")
	if title == "" || content == "" {
		return ctx.Redirect(http.StatusFound, "/blog/posts/create?message=标题和内容不能为空")
	}
	if len(title) > 40 {
		return ctx.Redirect(http.StatusFound, "/blog/posts/create?message=标题不能大于40个字符")
	}

	post := &models.Post{}
	err := post.CreatePost(c.db.ORM, c.user.ID, title, content, c.scRule.Rule.Post)
	if err != nil {
		return ctx.Redirect(http.StatusFound, "/blog/posts/create?message="+err.Error())
	}

	return ctx.Redirect(http.StatusFound, "/blog/posts")
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
	c.db.ORM.Model(post).Related(&comments)
	//c.db.ORM.Where("post_id = ?", postID).Find(&comments)
	author := new(models.User)
	c.db.ORM.Model(post).Related(author, "author_id")
	data := map[string]interface{}{
		"headTitle": c.title,
		"post":      post,
		"comments":  comments,
		"avatar":    functions.GetAvatarURL(c.user.Avatar),
		"author":    author,
	}
	post.CreatedAt.Format("")

	return ctx.Render(http.StatusOK, "posts/detail", data)
}

// UpdatePost CreatePost
func (c *PostController) UpdatePost(ctx gopress.Context) error {
	postIDStr := ctx.FormValue("post_id")
	postID, _ := strconv.ParseUint(postIDStr, 10, 64)
	post := &models.Post{}
	if c.db.ORM.First(post, postID).RecordNotFound() {
		return ctx.Redirect(http.StatusFound, "/blog/posts?message=文章不存在")
	}

	title := ctx.FormValue("title")
	content := ctx.FormValue("content")
	if title == "" || content == "" {
		return ctx.Redirect(http.StatusFound, "/blog/posts/update?message=标题和内容不能为空")
	}
	if len(title) > 40 {
		return ctx.Redirect(http.StatusFound, "/blog/posts/update?message=标题不能大于40个字符")
	}

	err := c.db.ORM.Save(post).Error
	if err != nil {
		return ctx.Redirect(http.StatusFound, "/blog/posts/update?message="+err.Error())
	}

	return ctx.Redirect(http.StatusFound, "/blog/posts/"+postIDStr)
}

// UpdatePage show create page
func (c *PostController) UpdatePage(ctx gopress.Context) error {
	postIDStr := ctx.Param("id")
	fmt.Println("")
	fmt.Println(postIDStr)
	fmt.Println("")
	postID, _ := strconv.ParseUint(postIDStr, 10, 64)
	post := &models.Post{}
	if c.db.ORM.First(post, postID).RecordNotFound() {
		return ctx.Redirect(http.StatusFound, "/assets/404.html")
	}

	data := map[string]interface{}{
		"headTitle": c.title,
		"avatar":    functions.GetAvatarURL(c.user.Avatar),
		"message":   ctx.QueryParam("message"),
		"post":      post,
	}

	return ctx.Render(http.StatusOK, "posts/update", data)
}

// MyPosts my posts
func (c *PostController) MyPosts(ctx gopress.Context) error {
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
	pl, err := p.MyPosts(c.db.ORM, pageIndex, limit, orderBy, c.user.ID)
	if err != nil {
		return ctx.Redirect(http.StatusFound, "/assets/404.html")
	}

	data := map[string]interface{}{
		"headTitle":    c.title,
		"avatar":       functions.GetAvatarURL(c.user.Avatar),
		"posts":        pl.Posts,
		"pagerContent": functions.GeneratePager(pl.Page, pl.Total, pl.Limit, pl.OrderBy, "/blog/posts", nil),
		"getAuthorName": func(id uint) string {
			return ""
		},
	}

	return ctx.Render(http.StatusOK, "posts/myposts", data)
}
