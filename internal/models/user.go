package models

import "time"

type User struct {
	ID               string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email            string    `gorm:"uniqueIndex;not null" json:"email"`
	Username         string    `gorm:"uniqueIndex" json:"username"`
	PasswordHash     string    `gorm:"not null" json:"-"`
	IsVerified       bool      `gorm:"default:true" json:"is_verified"`
	VerficationToken string    `gorm:"unique"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
