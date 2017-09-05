package models

import (
	"github.com/jinzhu/gorm"
)

type PostList struct {
	PageCount int
	Page      int
	Limit     int
	Total     int
	Posts     []*Post
	OrderBy   string
}

// Post table
type Post struct {
	gorm.Model
	Tittle        string `gorm:"not null;type:varchar(40);"`
	Content       string `gorm:"not null;type:text"`
	Comments      []Comment
	CommentNumber int `gorm:not null;default 0;index`
	AuthorID      int `gorm:"index"`
}

// TableName provide tabel naem to gorm
func (p *Post) TableName() string {
	return "posts"
}

// ListPosts list post
func (p *Post) ListPosts(ORM *gorm.DB, page, limit int, orderBy string) (*PostList, error) {
	if page < 1 {
		page = 1
	}
	start := limit * (page - 1)
	posts := []*Post{}

	if err := ORM.Offset(start).Limit(limit).Order(orderBy).Find(posts).Error; err != nil {
		return nil, DBError{"没有找到文章", DBReadError, err}
	}
	var total int
	ORM.Model(p).Count(&total)
	count := len(posts)

	pl := &PostList{}
	pl.Page = page
	pl.PageCount = count
	pl.Total = total
	pl.Posts = posts
	pl.OrderBy = orderBy

	return pl, nil
}
