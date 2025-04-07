package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type JWT struct {
	SecretKey string
}

type JWTData struct {
	Email string
}

func NewJWT(secretKey string) *JWT {
	return &JWT{
		SecretKey: secretKey,
	}
}

func (j *JWT) Creat(data *JWTData) (accessToken string, refreshToken string, error error) {
	accessClaims := jwt.MapClaims{
		"email": data.Email,
		"exp":   time.Now().Add(15 * time.Minute).Unix(),
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	accessToken, err := access.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", "", err
	}

	refreshToken = uuid.NewString()

	return accessToken, refreshToken, nil
}

//func (j *JWT) Parse(token string) (*JWTData, error) {
//	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
//		return []byte(j.SecretKey), nil
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	email := jwtToken.Claims.(jwt.MapClaims)["email"].(string)
//
//	return &JWTData{Email: email}, nil
//}
