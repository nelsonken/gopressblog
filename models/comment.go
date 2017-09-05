package models

import (
	"github.com/jinzhu/gorm"
)

// Comment table
type Comment struct {
	gorm.Model
	PostID   int `gorm:"index"`
	Content  string `gorm:"not null;type:varchar(200)"`
	AuthorID int `gorm:"index"`
}

// TableName provide tabel naem to gorm
func (u *Comment) TableName() string {
	return "comments"
}
