package services

import (
	"blog/models"
	"strconv"

	"github.com/fpay/gopress"
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

// GetCurrentUser get current user
func (s *UserService) GetCurrentUser(c gopress.Context) {
	s.User = &models.User{}
	cookie, _ := c.Cookie("uid")
	if cookie.Value == "" {
		return
	}
	container := gopress.AppFromContext(c).Services
	dbs := container.Get(DBServerName).(*DBService)
	uid, _ := strconv.Atoi(cookie.Value)

	dbs.ORM.First(s.User, uid)
}
