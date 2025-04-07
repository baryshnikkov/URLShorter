package auth

import (
	"URLShorter/internal/session"
	"URLShorter/pkg/req"
	"URLShorter/pkg/res"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net"
	"net/http"
)

type handler struct {
	Service        *Service
	SessionService *session.Service
}

type HandlerDeps struct {
	Service        *Service
	SessionService *session.Service
}

func NewHandler(router *chi.Mux, deps *HandlerDeps) {
	handler := &handler{
		Service:        deps.Service,
		SessionService: deps.SessionService,
	}
	router.HandleFunc("POST /auth/register", handler.Register())
	router.HandleFunc("POST /auth/login", handler.Login())
}

func (h *handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "auth.handler.Register"

		payload, err := req.JSON[RegisterReq](&w, r)
		if err != nil {
			return
		}

		registeredUser, err := h.Service.Register(payload.Email, payload.Login, payload.Password, payload.FirstName, payload.LastName)
		if err != nil {
			zap.L().Error("Error registering user",
				zap.String("op", op),
				zap.Error(err))
			res.JSON(w, err.Error(), http.StatusBadRequest)
			return
		}

		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip, _, _ = net.SplitHostPort(r.RemoteAddr)
		}
		userAgent := r.UserAgent()
		accessToken, refreshToken, err := h.SessionService.Save(registeredUser.ID, ip, userAgent, registeredUser.Email)
		if err != nil {
			zap.L().Error("Error registering session",
				zap.String("op", op),
				zap.Error(err))
			res.JSON(w, err.Error(), http.StatusInternalServerError)
		}

		registerRes := &RegisterRes{
			Email:        registeredUser.Email,
			Login:        registeredUser.Login,
			FirstName:    registeredUser.FirstName,
			LastName:     registeredUser.LastName,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}

		res.JSON(w, registerRes, http.StatusOK)
	}
}

func (h *handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "auth.handler.Login"

		payload, err := req.JSON[LoginReq](&w, r)
		if err != nil {
			return
		}

		user, err := h.Service.Login(payload.Email, payload.Password)
		if err != nil {
			zap.L().Error("Error login user",
				zap.String("op", op),
				zap.Error(err))
			res.JSON(w, err.Error(), http.StatusBadRequest)
			return
		}

		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip, _, _ = net.SplitHostPort(r.RemoteAddr)
		}
		userAgent := r.UserAgent()
		accessToken, refreshToken, err := h.SessionService.Save(user.ID, ip, userAgent, user.Email)
		if err != nil {
			zap.L().Error("Error registering session",
				zap.String("op", op),
				zap.Error(err))
			res.JSON(w, err.Error(), http.StatusInternalServerError)
		}

		loginRes := &LoginRes{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			Email:        user.Email,
			Login:        user.Login,
			FirstName:    user.FirstName,
			LastName:     user.LastName,
		}

		res.JSON(w, loginRes, http.StatusOK)
	}
}
