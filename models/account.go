package models

import (
	"github.com/jinzhu/gorm"
)

// Account table
type Account struct {
	gorm.Model
	OwnerID     int `gorm"index"`
	Total       float64
	TodayIncome float64
}
