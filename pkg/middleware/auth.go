package middleware

import (
	"URLShorter/configs"
	"URLShorter/internal/user"
	"URLShorter/pkg/jwt"
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type key string

const (
	ContextEmailKey  key = "email"
	ContextUserIdKey key = "user_id"
)

func IsAuthed(next http.Handler, config *configs.AppConfig, userRepository *user.Repository) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			writeUnauthorized(w, "unauthorized", errors.New("unauthorized"))
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		data, err := jwt.NewJWT(config.Auth.SecretKey).Parse(token)
		if err != nil {
			writeUnauthorized(w, "Invalid or expired token", err)
			return
		}

		user, err := userRepository.GetByEmail(data.Email)
		if err != nil || user == nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		ctx = context.WithValue(ctx, ContextUserIdKey, data.UserID)
		req := r.WithContext(ctx)

		next.ServeHTTP(w, req)
	})

}

func writeUnauthorized(w http.ResponseWriter, msg string, err error) {
	zap.L().Error(msg)
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(err.Error()))
}
