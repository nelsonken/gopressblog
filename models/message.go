package models

import (
	"github.com/jinzhu/gorm"
)

const (
	MessageTypeReciveComment      = 1
	MessageTypeReciveCommentReply = 2
	MessageTypeReciveGold         = 3
	MessageTypeSystem             = 4
)

// Message table
type Message struct {
	gorm.Model
	FromUserID  int
	ToUserID    int `gorm:"index"`
	Content     string
	MessageType uint
}
