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
	SystemUID = 0
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
func (m *Message) PutMessage(orm *gorm.DB, from, to uint, title, content string, messageType uint) error {
	m.FromUserID = from
	m.ToUserID = to
	m.Title = title
	m.Content = content
	m.MessageType = messageType
	m.Readed = false

	return orm.Save(m).Error
}

// ListMessages list some one's message, return total of the message
func (m *Message) ListMessages(orm *gorm.DB, userID uint, msgs *[]*Message, limit, page int, sortBy string) int {
	orm.Where("to_user_id = ?", userID).Order(sortBy).Offset(limit * (page - 1)).Limit(limit).Find(msgs)
	var total int
	orm.Model(m).Where("to_user_id = ?", userID).Count(&total)

	return total
}

// ReadMessage read a message
func (m *Message) ReadMessage(orm *gorm.DB, msgID uint) error {
	m.ID = msgID
	err := orm.Model(m).Update(map[string]interface{}{
		"readed": 1,
	}).Error

	return err
}

// ReadAll read all
func (m *Message) ReadAll(orm *gorm.DB, userID uint) error {
	return orm.Model(m).Where("to_user_id = ? AND readed = ?", userID, 0).UpdateColumn("readed", 1).Error
}

// DeleteOne remove oneMessage
func (m *Message) DeleteOne(orm *gorm.DB, msgID uint) error {
	m.ID = msgID
	return orm.Delete(m).Error
}
