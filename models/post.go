package models

import (
	"github.com/jinzhu/gorm"
)

// PostList paging query result
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
	Title         string `gorm:"not null;type:varchar(40);"`
	Content       string `gorm:"not null;type:text"`
	Comments      []Comment
	CommentNumber int  `gorm:"not null;default 0;index"`
	AuthorID      uint `gorm:"index"`
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
	if err := ORM.Offset(start).Limit(limit).Order(orderBy).Find(&posts).Error; err != nil {
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
	pl.Limit = limit

	return pl, nil
}

// CreatePost 新建文章
func (p *Post) CreatePost(ORM *gorm.DB, authorID uint, title, content string, score float64) error {
	p.Title = title
	p.Content = content
	p.AuthorID = authorID
	p.CommentNumber = 0

	tx := ORM.Begin()
	if err := tx.Save(p).Error; err != nil {
		tx.Rollback()
		return DBError{"保存失败", DBWriteError, err}
	}

	u := &Account{}

	if tx.First(u, "owner_id = ?", authorID).RecordNotFound() {
		u.OwnerID = authorID
		u.TodayIncome = score
		u.Total = score
		if err := tx.Save(u).Error; err != nil {
			tx.Rollback()
			return DBError{"保存失败", DBWriteError, err}
		}
	} else {
		u.Total += score
		u.TodayIncome += score
		if err := tx.Model(u).Updates(Account{Total: u.Total, TodayIncome: u.TodayIncome}).Error; err != nil {
			tx.Rollback()
			return DBError{"保存失败", DBWriteError, err}
		}
	}
	tx.Commit()

	return nil
}
