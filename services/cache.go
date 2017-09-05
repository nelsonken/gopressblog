package services

import (
	"github.com/fpay/gopress"
)

const (
	// CacheServiceName is the identity of cache service
	CacheServiceName = "cache"
)

// CacheService type
type CacheService struct {
	// Uncomment this line if this service has dependence on other services in the container
	// c *gopress.Container
}

// NewCacheService returns instance of cache service
func NewCacheService() *CacheService {
	return new(CacheService)
}

// ServiceName is used to implements gopress.Service
func (s *CacheService) ServiceName() string {
	return CacheServiceName
}

// RegisterContainer is used to implements gopress.Service
func (s *CacheService) RegisterContainer(c *gopress.Container) {
	// Uncomment this line if this service has dependence on other services in the container
	// s.c = c
}

func (s *CacheService) SampleMethod() {
}
