package session

import (
	"time"

	"URLShorter/internal/user"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model

	UserID           uint      `gorm:"index;not null"`
	RefreshTokenHash string    `gorm:"size:255;uniqueIndex;not null"`
	IP               string    `gorm:"size:45"`
	UserAgent        string    `gorm:"size:255"`
	ExpiresAt        time.Time `gorm:"not null"`

	User user.User `gorm:"constraint:OnDelete:CASCADE"`
}
