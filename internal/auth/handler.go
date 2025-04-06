package auth

import (
	"URLShorter/pkg/req"
	"URLShorter/pkg/res"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type handler struct {
	Service *Service
}

type HandlerDeps struct {
	Service *Service
}

func NewHandler(router *chi.Mux, deps *HandlerDeps) {
	handler := &handler{
		Service: deps.Service,
	}
	router.HandleFunc("POST /auth/register", handler.Register())
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

		registerRes := &RegisterRes{
			Email:     registeredUser.Email,
			Login:     registeredUser.Login,
			FirstName: registeredUser.FirstName,
			LastName:  registeredUser.LastName,
		}

		res.JSON(w, registerRes, http.StatusOK)
	}
}
