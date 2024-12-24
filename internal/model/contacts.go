package model

import "time"

type Contact struct {
	ID        int64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"not null" json:"name"`
	Email     string `gorm:"not null" json:"email"`
	Phone     string `gorm:"not null" json:"phone"`
	GedungId  int64  `gorm:"not null" json:"gedung_id"`
	KotaId    int64  `gorm:"not null" json:"kota_id"`
	Role      string `gorm:"not null;default:CUSTOMER_SERVICE" json:"role"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (Contact) TableName() string {
	return "contacts"
}
