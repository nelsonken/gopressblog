package services

import (
	"github.com/fpay/gopress"
	"gopkg.in/go-playground/validator.v9"
)

const (
	// ValidatorServiceName is the identity of vlidator service
	ValidatorServiceName = "vlidator"
)

// Validator type
type Validator struct {
	// Uncomment this line if this service has dependence on other services in the container
	// c *gopress.Container
	V *validator.Validate
}

// NewValidatorService returns instance of vlidator service
func NewValidatorService() *Validator {
	v := &Validator{}
	v.V = validator.New()

	return v
}

// ServiceName is used to implements gopress.Service
func (s *Validator) ServiceName() string {
	return ValidatorServiceName
}

// RegisterContainer is used to implements gopress.Service
func (s *Validator) RegisterContainer(c *gopress.Container) {
	// Uncomment this line if this service has dependence on other services in the container
	// s.c = c
}

// Validate validate a struct data
func (s *Validator) Validate(i interface{}) error {
	return s.V.Struct(i)
}
