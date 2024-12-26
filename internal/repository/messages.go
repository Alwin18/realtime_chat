package repository

import (
	"github.com/websoket-chat/internal/model"
	"gorm.io/gorm"
)

type ChatMessageRepository struct {
	db *gorm.DB
}

func NewChatMessageRepository(db *gorm.DB) *ChatMessageRepository {
	return &ChatMessageRepository{db: db.Debug()}
}

func (repo *ChatMessageRepository) SaveMessage(message model.Message) error {
	return repo.db.Create(&message).Error
}

func (repo *ChatMessageRepository) GetChatHistory(senderID string, receiverID string) ([]model.HistoryChat, error) {
	var messages []model.HistoryChat
	err := repo.db.Table("messages m").
		Select(`
			m.id AS message_id,			
			m.content,
			m.attachment_url,
			m.sent_at,
			m.is_read
	`).
		Where(`
			(m.sender_id = ? AND m.receiver_id = ?) 
    		OR 
    		(m.sender_id = ? AND m.receiver_id = ?)
		`, senderID, receiverID, receiverID, senderID).Order("m.sent_at desc").Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}
