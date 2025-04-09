package link

import (
	"URLShorter/internal/user"
	"encoding/base64"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Link struct {
	gorm.Model

	UserID uint   `gorm:"index;not null"`
	Url    string `gorm:"size:255;not null"`
	Hash   string `gorm:"uniqueIndex"`

	User user.User `gorm:"constraint:OnDelete:CASCADE"`
}

func NewLink(userId uint, url string) *Link {
	link := &Link{
		UserID: userId,
		Url:    url,
	}

	link.setHash()

	return link
}

func (l *Link) setHash() {
	u := uuid.New()
	l.Hash = base64.URLEncoding.EncodeToString(u[:])[:8]
}
