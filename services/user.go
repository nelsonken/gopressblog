package services

import (
	"blog/models"
	"github.com/fpay/gopress"
	"net/http"
	"strconv"
)

const (
	// UserServiceName is the identity of user service
	UserServiceName = "user"
)

// UserService type
type UserService struct {
	User *models.User
}

// NewUserService returns instance of user service
func NewUserService() *UserService {

	return new(UserService)
}

// ServiceName is used to implements gopress.Service
func (s *UserService) ServiceName() string {
	return UserServiceName
}

// RegisterContainer is used to implements gopress.Service
func (s *UserService) RegisterContainer(c *gopress.Container) {
	// Uncomment this line if this service has dependence on other services in the container
	// s.c = c
}

func (s *UserService) GetCurrentUser(c gopress.Context) error {
	cookie, err := c.Cookie("uid")
	if err != nil {
		return c.Redirect(http.StatusFound, "/login")
	}

	if cookie.Value == "" {
		return c.Redirect(http.StatusFound, "/login")
	}

	container := gopress.AppFromContext(c).Services

	dbs := container.Get(DBServerName).(*DBService)
	uid, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return c.Redirect(http.StatusFound, "/login")
	}
	if dbs.ORM.First(s.User, uid).RecordNotFound() {
		return c.Redirect(http.StatusFound, "/login")
	}

	return nil
}
