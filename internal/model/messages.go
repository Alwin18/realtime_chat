package model

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	SenderID      uuid.UUID  `gorm:"type:uuid;not null" json:"sender_id"`
	ReceiverID    uuid.UUID  `gorm:"type:uuid;not null" json:"receiver_id"`
	Content       *string    `gorm:"type:text" json:"content"`
	AttachmentURL *string    `gorm:"type:varchar(255)" json:"attachment_url"`
	SentAt        time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"sent_at"`
	IsRead        bool       `gorm:"default:false" json:"is_read"`
	ReadAt        *time.Time `gorm:"type:timestamp" json:"read_at"`
	Sender        User       `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE" json:"sender"`
	Receiver      User       `gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE" json:"receiver"`
}

func (Message) TableName() string {
	return "messages"
}

type HistoryChat struct {
	MessageID     int64     `json:"message_id"`
	Content       string    `json:"content"`
	AttachmentURL string    `json:"attachment_url"`
	SentAt        time.Time `json:"sent_at"`
	IsRead        bool      `json:"is_read"`
}
