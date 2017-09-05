package services

import (
	"github.com/fpay/gopress"
	"gopkg.in/go-playground/validator.v9"
)

const (
	// VlidatorServiceName is the identity of vlidator service
	VlidatorServiceName = "vlidator"
)

// VlidatorService type
type VlidatorService struct {
	// Uncomment this line if this service has dependence on other services in the container
	// c *gopress.Container
	V *validator.Validate
}

// NewVlidatorService returns instance of vlidator service
func NewVlidatorService() *VlidatorService {
	v := &VlidatorService{}
	v.V = validator.New()

	return v
}

// ServiceName is used to implements gopress.Service
func (s *VlidatorService) ServiceName() string {
	return VlidatorServiceName
}

// RegisterContainer is used to implements gopress.Service
func (s *VlidatorService) RegisterContainer(c *gopress.Container) {
	// Uncomment this line if this service has dependence on other services in the container
	// s.c = c
}

func (s *VlidatorService) Validate(i interface{}) error {
	return s.V.Struct(i)
}
