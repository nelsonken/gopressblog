package models

import (
	"github.com/jinzhu/gorm"
)

const (
	PostIncome = 10
	CommentIncome = 5
)

// Account table
type Account struct {
	gorm.Model
	OwnerID     int `gorm"index"`
	Total       float64
	TodayIncome float64
}
