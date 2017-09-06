package models

import (
	"github.com/jinzhu/gorm"
)

const (
	// MessageTypeReciveComment got a comment
	MessageTypeReciveComment = 1
	// MessageTypeReciveCommentReply got a reply
	MessageTypeReciveCommentReply = 2
	// MessageTypeGotScore get a score
	MessageTypeGotScore = 3
	// MessageTypeSystem  system msg
	MessageTypeSystem = 4
	// SystemUID systemUID
	SystemUID = 1
)

// Message table
type Message struct {
	gorm.Model
	FromUserID  uint
	ToUserID    uint   `gorm:"index"`
	Title       string `gorm:"not null;type:varchar(50)"`
	Content     string
	MessageType uint
	Readed      bool
}

// PutMessage put message to
func (m *Message) PutMessage(ORM *gorm.DB, from, to uint, title, content string, messageType uint) error {
	m.FromUserID = from
	m.ToUserID = to
	m.Title = title
	m.Content = content
	m.MessageType = messageType
	m.Readed = false

	return ORM.Save(m).Error
}
