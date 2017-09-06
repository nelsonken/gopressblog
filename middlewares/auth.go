package middlewares

import (
	"blog/models"
	"blog/services"
	"github.com/fpay/gopress"
	"net/http"
	"strconv"
	"time"
)

// NewAuthMiddleware returns auth middleware.
// check cookies
func NewAuthMiddleware(user *models.User) gopress.MiddlewareFunc {
	return func(next gopress.HandlerFunc) gopress.HandlerFunc {
		return func(c gopress.Context) error {
			panic("hello")
			cookie, err := c.Cookie("uid")
			if err != nil {
				return c.Redirect(http.StatusFound, "/login")
			}

			if cookie.Value == "" {
				dropCookie(c, cookie)
				return c.Redirect(http.StatusFound, "/login")
			}

			container := gopress.AppFromContext(c).Services

			dbs := container.Get(services.DBServerName).(*services.DBService)
			uid, err := strconv.Atoi(cookie.Value)
			if err != nil {
				dropCookie(c, cookie)
				return c.Redirect(http.StatusFound, "/login")
			}
			if dbs.ORM.First(user, uid).RecordNotFound() {
				dropCookie(c, cookie)
				return c.Redirect(http.StatusFound, "/login")
			}
			us := services.NewUserService()
			us.User = user
			return next(c)
		}
	}
}

func dropCookie(c gopress.Context, cookie *http.Cookie) {
	cookie.Expires = time.Now().Add(-1 * time.Second)
	c.SetCookie(cookie)
}
