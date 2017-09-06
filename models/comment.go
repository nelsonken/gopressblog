package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

// Comment table
type Comment struct {
	gorm.Model
	PostID      uint   `gorm:"not null; default 0;index"`
	TartgetID   uint   `gorm:"not null; default 0;index"`
	Content     string `gorm:"not null;type:varchar(200)"`
	AuthorID    uint   `gorm:"index"`
	CommentType uint
}

// TableName provide tabel naem to gorm
func (c *Comment) TableName() string {
	return "comments"
}

// CommetnPost create a comment for post
func (c *Comment) CommentPost(orm *gorm.DB, postID, authorID, mentionUserID uint, content string, score float64) error {
	c.AuthorID = authorID
	c.PostID = postID
	c.Content = content
	c.TartgetID = mentionUserID
	ta := orm.Begin()

	// 增加积分
	account := &Account{}
	if ta.Set("gorm:query_option", "FOR UPDATE").FirstOrCreate(account, "owner_id = ?", authorID).Error != nil {
		return DBError{"评论作者已不存在", DBReadError, nil}
	}
	account.TodayIncome += score
	account.Total = score
	ea := ta.Model(account).Update(&Account{Total: account.Total, TodayIncome: account.TodayIncome}).Error

	// 增加评论次数
	post := &Post{}
	if ta.Set("gorm:query_option", "FOR UPDATE").First(post, c.PostID).RecordNotFound() {
		return DBError{"文章已不存在", DBReadError, nil}
	}
	ep := ta.Model(post).Update(&Post{CommentNumber: post.CommentNumber + 1}).Error

	// 发送评论提醒
	var commenters []uint
	msg := &Message{}
	ta.Model(c).Pluck("author_id", &commenters)
	for _, v := range commenters {
		if mentionUserID == v {
			msgTitle := fmt.Sprintf("作者在文章《%s》中回复了了你的评论", post.Title)
			msg.PutMessage(ta, SystemUID, mentionUserID, msgTitle, content, MessageTypeSystem)
			break
		}
	}

	// 提醒作者被评论
	msgTitle := fmt.Sprintf("你的文章《%s》有一个新评论", post.Title)
	msg.PutMessage(ta, SystemUID, post.AuthorID, msgTitle, content, MessageTypeSystem)

	// 新建评论
	if ta.Save(c).Error != nil || ea != nil || ep != nil {
		ta.Rollback()
		return DBError{"存储评论失败", DBWriteError, nil}
	}

	// 检查错误
	errs := ta.GetErrors()
	if len(errs) > 0 {
		// 回滚
		ta.Rollback()
		return DBError{"存储评论失败", DBWriteError, nil}
	}
	// 提交
	ta.Commit()

	return nil
}
