package models

import (
	"time"
)

type UserRole struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Role      string    `gorm:"not null" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
