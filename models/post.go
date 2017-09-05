package models

import (
	"github.com/jinzhu/gorm"
)

// Post table
type Post struct {
	gorm.Model
	Tittle        string `gorm:"not null;type:varchar(40);"`
	Content       string `gorm:"not null;type:text"`
	Comments      []Comment
	CommentNumber int
	AuthorID      int `gorm:"index"`
}

// TableName provide tabel naem to gorm
func (u *Post) TableName() string {
	return "posts"
}
