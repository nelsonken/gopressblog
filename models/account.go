package models

import (
	"github.com/jinzhu/gorm"
)

// Account table
type Account struct {
	gorm.Model
	OwnerID     uint `gorm:"unique"`
	Total       float64
	TodayIncome float64
}
