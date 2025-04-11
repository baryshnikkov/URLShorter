package session

import (
	"URLShorter/pkg/res"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net"
	"net/http"
)

type Handler struct {
	Service *Service
}

type HandlerDeps struct {
	Service *Service
}

func NewHandler(router *chi.Mux, deps *HandlerDeps) {
	handler := &Handler{
		Service: deps.Service,
	}
	router.HandleFunc("POST /session/refresh/{token}", handler.UpdateRefreshToken())
}

func (h *Handler) UpdateRefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "session.Handler.UpdateRefreshToken"

		token := r.PathValue("token")
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip, _, _ = net.SplitHostPort(r.RemoteAddr)
		}
		userAgent := r.UserAgent()

		accessTokenNew, refreshTokenNew, err := h.Service.UpdateRefreshToken(token, ip, userAgent)
		if err != nil {
			zap.L().Error("Error update refresh token",
				zap.String("op", op),
				zap.String("token", token),
				zap.Error(err))
			res.JSON(w, "err", http.StatusBadRequest)
			return
		}

		response := &updateRefreshTokenRes{
			RefreshToken: refreshTokenNew,
			AccessToken:  accessTokenNew,
		}

		res.JSON(w, response, http.StatusOK)
	}
}
