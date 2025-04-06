package user

import (
	"URLShorter/pkg/di"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email             string  `gorm:"size:255;uniqueIndex;not null"`
	Login             string  `gorm:"size:50;uniqueIndex;not null"`
	PasswordHash      string  `gorm:"size:255;not null"`
	FirstName         string  `gorm:"size:100"`
	LastName          string  `gorm:"size:100"`
	Role              di.Role `gorm:"type:varchar(20);default:user;not null"`
	IsBanned          bool    `gorm:"default:false"`
	EmailVerified     bool    `gorm:"default:false"`
	VerificationToken string  `gorm:"size:255"`
}
