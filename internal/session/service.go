package session

import (
	"URLShorter/configs"
	"URLShorter/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service struct {
	Repository *Repository
	AppConfig  *configs.AppConfig
}

type ServiceDeps struct {
	Repository *Repository
	AppConfig  *configs.AppConfig
}

func NewService(deps *ServiceDeps) *Service {
	return &Service{
		Repository: deps.Repository,
		AppConfig:  deps.AppConfig,
	}
}

func (s *Service) Save(userID uint, ip string, userAgent string, email string) (accessToken string, refreshToken string, error error) {
	accessToken, refreshToken, err := jwt.NewJWT(s.AppConfig.Auth.SecretKey).Creat(&jwt.JWTData{Email: email})
	if err != nil {
		return "", "", err
	}

	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	_, err = s.Repository.Save(&Session{
		UserID:           userID,
		IP:               ip,
		UserAgent:        userAgent,
		ExpiresAt:        time.Now().Add(7 * 24 * time.Hour),
		RefreshTokenHash: string(refreshTokenHash),
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
