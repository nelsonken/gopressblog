package controllers

import (
	"net/http"
	"strconv"

	"github.com/fpay/gopress"
	"github.com/labstack/echo"

	"blog/functions"
	"blog/models"
	"blog/services"
	"fmt"
)

// PostController post controller
type PostController struct {
	db     *services.DBService
	title  string
	scRule *services.ScoreService
	group  *echo.Group
	es     *services.ElasticService
}

// NewPostController returns post controller instance.
func NewPostController(group *echo.Group) *PostController {
	c := new(PostController)
	c.group = group

	group.GET("/posts", c.ListPosts)
	group.GET("/posts/search", c.SearchPost)
	group.GET("/posts/create", c.CreatePage)
	group.GET("/posts/update/:id", c.UpdatePage)
	group.POST("/posts/update", c.UpdatePost)
	group.GET("/posts/:id", c.ViewPost)
	group.POST("/posts/create", c.CreatePost)
	group.GET("/myposts", c.MyPosts)

	return c
}

// RegisterRoutes registes routes to app
// It is used to implements gopress.Controller.
func (c *PostController) RegisterRoutes(app *gopress.App) {
	c.db = app.Services.Get(services.DBServerName).(*services.DBService)
	c.scRule = app.Services.Get(services.ScoreServiceName).(*services.ScoreService)
	c.es = app.Services.Get(services.ElasticServiceName).(*services.ElasticService)
	c.title = "BLOG-Article"
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
		orderBy = defaultSortBy
	}

	limit := 10
	p := &models.Post{}
	pl, err := p.ListPosts(c.db.ORM, pageIndex, limit, orderBy)
	if err != nil {
		return ctx.Redirect(http.StatusFound, notFoundURL)
	}

	data := map[string]interface{}{
		"headTitle":        c.title,
		"haveMessage":      ctx.Get("haveMessage"),
		"messageNum":       ctx.Get("messageNum"),
		"avatar":           functions.GetAvatarURL(getUser(ctx).Avatar),
		"posts":            pl.Posts,
		"pagerContent":     functions.GeneratePager(pl.Page, pl.Total, pl.Limit, pl.OrderBy, "/blog/posts", nil),
		"hotAuthors":       c.getHotAuthors(),
		"getUserName":      c.getUserName,
		"getUserAvatarURL": c.getUserAvatarURL,
		"friendlyTime":     functions.GetFriendlyTime,
	}

	return ctx.Render(http.StatusOK, "posts/list", data)
}

func (c *PostController) getHotAuthors() []*models.User {
	hotAuthorsID := []uint{}
	c.db.ORM.Model(&models.Account{}).Order("today_income desc").Limit(10).Pluck("owner_id", &hotAuthorsID)
	hotAuthors := []*models.User{}
	c.db.ORM.Select("id, name").Where("id in (?)", hotAuthorsID).Find(&hotAuthors)
	return hotAuthors
}

// CreatePage show create page
func (c *PostController) CreatePage(ctx gopress.Context) error {
	data := map[string]interface{}{
		"headTitle":   c.title,
		"avatar":      functions.GetAvatarURL(getUser(ctx).Avatar),
		"haveMessage": ctx.Get("haveMessage"),
		"messageNum":  ctx.Get("messageNum"),
		"message":     ctx.QueryParam("message"),
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
	err := post.CreatePost(c.db.ORM, getUser(ctx).ID, title, content, c.scRule.Rule.Post)
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
		return ctx.Redirect(http.StatusFound, notFoundURL)
	}

	post := &models.Post{}
	if c.db.ORM.First(post, postID).RecordNotFound() {
		return ctx.Redirect(http.StatusFound, notFoundURL)
	}

	comments := []*models.Comment{}
	c.db.ORM.Model(post).Order("created_at desc").Related(&comments)
	author := new(models.User)
	c.db.ORM.Model(post).Related(author, "author_id")
	commentator := []uint{}
	c.db.ORM.Model(&models.Comment{}).Where("post_id = ? ", post.ID).Pluck("Distinct(author_id)", &commentator)
	data := map[string]interface{}{
		"headTitle":        c.title,
		"post":             post,
		"comments":         comments,
		"haveMessage":      ctx.Get("haveMessage"),
		"messageNum":       ctx.Get("messageNum"),
		"avatar":           functions.GetAvatarURL(getUser(ctx).Avatar),
		"author":           author,
		"commentator":      commentator,
		"getUserName":      c.getUserName,
		"message":          ctx.QueryParam("message"),
		"getUserAvatarURL": c.getUserAvatarURL,
	}

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
		return ctx.Redirect(http.StatusFound, notFoundURL)
	}

	data := map[string]interface{}{
		"headTitle":   c.title,
		"avatar":      functions.GetAvatarURL(getUser(ctx).Avatar),
		"message":     ctx.QueryParam("message"),
		"haveMessage": ctx.Get("haveMessage"),
		"messageNum":  ctx.Get("messageNum"),
		"post":        post,
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
		orderBy = defaultSortBy
	}

	limit := 10
	p := &models.Post{}
	pl, err := p.MyPosts(c.db.ORM, pageIndex, limit, orderBy, getUser(ctx).ID)
	if err != nil {
		return ctx.Redirect(http.StatusFound, notFoundURL)
	}

	data := map[string]interface{}{
		"headTitle":    c.title,
		"haveMessage":  ctx.Get("haveMessage"),
		"messageNum":   ctx.Get("messageNum"),
		"avatar":       functions.GetAvatarURL(getUser(ctx).Avatar),
		"posts":        pl.Posts,
		"pagerContent": functions.GeneratePager(pl.Page, pl.Total, pl.Limit, pl.OrderBy, "/blog/posts", nil),
		"getUserName":  c.getUserName,
	}

	return ctx.Render(http.StatusOK, "posts/myposts", data)
}

func (c *PostController) getUserName(uID uint) string {
	u := &models.User{}
	c.db.ORM.Select("name").First(u, uID)
	return u.Name
}

// SearchPost search posts
func (c *PostController) SearchPost(ctx gopress.Context) error {
	keyword := ctx.QueryParam("keyword")
	pageS := ctx.QueryParam("page")
	page, _ := strconv.Atoi(pageS)
	if page < 1 {
		page++
	}

	limit := 10
	posts, err := c.es.SearchPosts(keyword, limit, page)

	data := map[string]interface{}{
		"headTitle":   c.title,
		"haveMessage": ctx.Get("haveMessage"),
		"messageNum":  ctx.Get("messageNum"),
		"avatar":      functions.GetAvatarURL(getUser(ctx).Avatar),
		"posts":       posts,
		"keyword":     keyword,
		"hotAuthors":  c.getHotAuthors(),
		"getUserName": c.getUserName,
	}
	if err != nil {
		data["message"] = err.Error()
	}

	return ctx.Render(http.StatusOK, "posts/search", data)
}

func (c *PostController) getUserAvatarURL(uID uint) string {
	u := &models.User{}
	c.db.ORM.Select("avatar").First(u, uID)

	return functions.GetAvatarURL(u.Avatar)
}
