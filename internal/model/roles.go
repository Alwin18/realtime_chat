package model

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Role      string    `gorm:"type:varchar(20);not null;check:role IN ('CUSTOMER_SERVICE','SISWA','ORANG_TUA')" json:"role"`
	CreatedAt time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	Users     []User    `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"users,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}
