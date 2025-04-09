package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

type JWT struct {
	SecretKey string
}

type JWTData struct {
	Email     string
	ExpiresAt int64
	UserID    uint
}

func NewJWT(secretKey string) *JWT {
	return &JWT{
		SecretKey: secretKey,
	}
}

func (j *JWT) Creat(data *JWTData) (accessToken string, refreshToken string, error error) {
	const op string = "jwt.Create"

	accessClaims := jwt.MapClaims{
		"email":   data.Email,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"user_id": int64(data.UserID),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	accessToken, err := access.SignedString([]byte(j.SecretKey))
	if err != nil {
		zap.L().Error("Error creating access token",
			zap.String("op", op),
			zap.Error(err))
		return "", "", err
	}

	refreshToken = uuid.NewString()

	return accessToken, refreshToken, nil
}

func (j *JWT) Parse(accessToken string) (*JWTData, error) {
	const op string = "jwt.Parse"

	jwtToken, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		zap.L().Error("Error parsing access token",
			zap.String("op", op),
			zap.Error(err))
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		zap.L().Error("Error parsing access token (claims)",
			zap.String("op", op),
			zap.Error(err))
		return nil, errors.New("invalid token")
	}

	rawUserID, ok := claims["user_id"]
	if !ok {
		zap.L().Error("user_id not found in token", zap.String("op", op))
		return nil, errors.New("user_id not found in token")
	}
	userIDFloat, ok := rawUserID.(float64)
	if !ok {
		zap.L().Error("invalid user_id type", zap.String("op", op))
		return nil, errors.New("invalid user_id type")
	}

	email, ok := claims["email"].(string)
	if !ok {
		zap.L().Error("Error parsing email claim",
			zap.String("op", op),
			zap.Error(err))
		return nil, errors.New("email not found in token")
	}

	expFloat, ok := claims["exp"].(float64)
	if !ok {
		zap.L().Error("Error parsing exp claim",
			zap.String("op", op),
			zap.Error(err))
		return nil, errors.New("exp not found in token")
	}

	exp := time.Unix(int64(expFloat), 0)
	if time.Now().After(exp) {
		return nil, errors.New("token expired")
	}

	return &JWTData{
		Email:     email,
		ExpiresAt: exp.Unix(),
		UserID:    uint(userIDFloat),
	}, nil
}
