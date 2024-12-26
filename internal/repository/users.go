package repository

import (
	"github.com/websoket-chat/internal/model"
	"gorm.io/gorm"
)

type ContactRepository struct {
	db *gorm.DB
}

func NewContactRepository(db *gorm.DB) *ContactRepository {
	return &ContactRepository{db: db.Debug()}
}

func (repo *ContactRepository) GetContactByCakupan(gedungId int64, role string) ([]model.User, error) {
	var contacts []model.User
	query := repo.db.Model(&model.User{})

	if gedungId != 0 {
		query = query.Where("gedung_id = ?", gedungId)
	}

	if role == "ORANG_TUA" {
		query = query.Where("role = ?", role)
	}

	if err := query.Order("name asc").Find(&contacts).Error; err != nil {
		return nil, err
	}

	return contacts, nil
}
