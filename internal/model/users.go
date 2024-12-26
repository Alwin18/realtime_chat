package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name             string     `gorm:"type:varchar(100);not null" json:"name"`
	Email            *string    `gorm:"type:varchar(255)" json:"email,omitempty"`
	PhoneNumber      string     `gorm:"type:varchar(14);not null;unique" json:"phone_number"`
	AvatarURL        *string    `gorm:"type:varchar(255)" json:"avatar_url,omitempty"`
	KotaID           int64      `gorm:"type:bigint;not null;default:1" json:"kota_id"`
	GedungID         int64      `gorm:"type:bigint;not null;default:1" json:"gedung_id"`
	IsOnline         bool       `gorm:"default:false" json:"is_online"`
	RoleID           uuid.UUID  `gorm:"type:uuid;not null" json:"role_id"`
	Role             Role       `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"role,omitempty"`
	CreatedAt        time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        *time.Time `gorm:"type:timestamp" json:"updated_at,omitempty"`
	SentMessages     []Message  `gorm:"foreignKey:SenderID;constraint:OnDelete:CASCADE" json:"sent_messages,omitempty"`
	ReceivedMessages []Message  `gorm:"foreignKey:ReceiverID;constraint:OnDelete:CASCADE" json:"received_messages,omitempty"`
}

func (User) TableName() string {
	return "users"
}
