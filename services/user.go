package services

import (
	"github.com/fpay/gopress"
	"blog/models"
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