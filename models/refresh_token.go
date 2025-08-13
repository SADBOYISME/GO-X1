package models

import (
	"time"

	"gorm.io/gorm"
)

// RefreshToken stores a refresh token and its state
type RefreshToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Token     string    `json:"token" gorm:"unique;not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	Revoked   bool      `json:"revoked" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// RevokeRefreshToken marks a refresh token as revoked
func RevokeRefreshToken(db *gorm.DB, token string) error {
	return db.Model(&RefreshToken{}).Where("token = ?", token).Update("revoked", true).Error
}
