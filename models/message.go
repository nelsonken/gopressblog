package models

import (
	"github.com/jinzhu/gorm"
)

const (
	MessageTypeReciveComment      = 1
	MessageTypeReciveCommentReply = 2
	MessageTypeReciveGold         = 3
	MessageTypeSystem             = 4
	SystemFrom                    = 1
)

// Message table
type Message struct {
	gorm.Model
	FromUserID  int
	ToUserID    int    `gorm:"index"`
	Title       string `gorm:"not null;type:varchar(50)"`
	Content     string
	MessageType uint
	Readed      bool
}

// PutMessage put message to
func (m *Message) PutMessage(ORM *gorm.DB, from, to int, title, content string, messageType uint) error {
	m.FromUserID = from
	m.ToUserID = to
	m.Title = title
	m.Content = content
	m.MessageType = messageType
	m.Readed = false

	return ORM.Save(m).Error
}
