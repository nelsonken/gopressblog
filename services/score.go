package services

import (
	"github.com/fpay/gopress"
)

const (
	// ScoreServiceName is the identity of score service
	ScoreServiceName = "score"
)

// ScoreService type
type ScoreService struct {
	// Uncomment this line if this service has dependence on other services in the container
	Rule *ScoreRule
}

// ScoreRule 积分奖励规则
type ScoreRule struct {
	Post    float64 `yaml:"post"`
	Comment float64 `yaml:"comment"`
}

// NewScoreService returns instance of score service
func NewScoreService(s *ScoreRule) *ScoreService {
	return &ScoreService{s}
}

// ServiceName is used to implements gopress.Service
func (s *ScoreService) ServiceName() string {
	return ScoreServiceName
}

// RegisterContainer is used to implements gopress.Service
func (s *ScoreService) RegisterContainer(c *gopress.Container) {
	// Uncomment this line if this service has dependence on other services in the container
	// s.c = c
}
