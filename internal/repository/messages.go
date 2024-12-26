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

func (repo *ChatMessageRepository) GetChatHistory(senderID string, receiverID string) ([]model.Message, error) {
	var messages []model.Message
	err := repo.db.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		senderID, receiverID, receiverID, senderID,
	).Order("timestamp desc").Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}
