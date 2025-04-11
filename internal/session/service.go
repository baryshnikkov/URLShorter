package session

import (
	"URLShorter/configs"
	"URLShorter/internal/user"
	"URLShorter/pkg/jwt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	Repository     *Repository
	UserRepository *user.Repository
	AppConfig      *configs.AppConfig
}

type ServiceDeps struct {
	Repository     *Repository
	UserRepository *user.Repository
	AppConfig      *configs.AppConfig
}

func NewService(deps *ServiceDeps) *Service {
	return &Service{
		Repository:     deps.Repository,
		UserRepository: deps.UserRepository,
		AppConfig:      deps.AppConfig,
	}
}

func (s *Service) Save(userID uint, ip string, userAgent string, email string) (accessToken string, refreshToken string, error error) {
	const op string = "session.Save"

	accessToken, refreshToken, err := jwt.NewJWT(s.AppConfig.Auth.SecretKey).Create(&jwt.JWTData{
		Email:  email,
		UserID: userID,
	})
	if err != nil {
		zap.L().Error("Error create new JWT",
			zap.String("op", op),
			zap.Error(err))
		return "", "", err
	}

	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error("Error compare JWT tokens",
			zap.String("op", op),
			zap.Error(err))
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
		zap.L().Error("Error save session",
			zap.String("op", op),
			zap.Error(err))
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *Service) UpdateRefreshToken(refreshToken string, ip string, userAgent string) (accessTokenNew string, refreshTokenNew string, err error) {
	const op = "session.service.UpdateRefreshToken"

	userIdStr := strings.Split(refreshToken, "&")[0]
	userIdUint64, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		zap.L().Error("Error convert refresh token to uint",
			zap.String("op", op),
			zap.Error(err))
		return "", "", err
	}
	userId := uint(userIdUint64)

	sessions, err := s.Repository.FindAllByUserId(userId)
	if err != nil {
		zap.L().Error("Error get all sessions by userId",
			zap.String("op", op),
			zap.Error(err))
		return "", "", err
	}

	var coincidentSession *Session
	for i, session := range sessions {
		err = bcrypt.CompareHashAndPassword([]byte(session.RefreshTokenHash), []byte(refreshToken))
		if err != nil {
			coincidentSession = sessions[i]
		}
	}
	if coincidentSession == nil {
		zap.L().Error("can not compare refresh token",
			zap.String("op", op),
			zap.Error(err))
		return "", "", err
	}

	err = s.Repository.DeleteByRefreshTokenHash(coincidentSession.RefreshTokenHash)
	if err != nil {
		zap.L().Error("Error delete sessions by refresh token hash",
			zap.String("op", op),
			zap.Error(err))
		return "", "", err
	}

	email, err := s.UserRepository.GetEmailById(userId)

	accessTokenNew, refreshTokenNew, err = s.Save(userId, ip, userAgent, email)
	if err != nil {
		zap.L().Error("Error create new session",
			zap.String("op", op),
			zap.Error(err))
		return "", "", err
	}

	return accessTokenNew, refreshTokenNew, nil
}
