package model

import "time"

type ChatMessage struct {
	ID         int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	SenderID   int64     `gorm:"not null" json:"sender_id"`
	ReceiverID int64     `gorm:"not null" json:"receiver_id"`
	Content    string    `gorm:"not null" json:"content"`
	IsRead     bool      `gorm:"not null" json:"is_read"`
	Timestamp  time.Time `gorm:"autoCreateTime" json:"timestamp"`
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}
