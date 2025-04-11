package link

import (
	"URLShorter/configs"
	"URLShorter/internal/user"
	"URLShorter/pkg/middleware"
	"URLShorter/pkg/req"
	"URLShorter/pkg/res"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	Repository *Repository
}

type HandlerDeps struct {
	Repository     *Repository
	UserRepository *user.Repository
	AppConfig      *configs.AppConfig
}

func NewHandler(router *chi.Mux, deps *HandlerDeps) {
	handler := &Handler{
		Repository: deps.Repository,
	}

	router.Get("/{hash}", handler.goTo())

	router.Group(func(rg chi.Router) {
		rg.Use(func(next http.Handler) http.Handler {
			return middleware.IsAuthed(next, deps.AppConfig, deps.UserRepository)
		})

		rg.Post("/link", handler.create())
	})
}

func (h *Handler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "link.handler.create()"

		userID, ok := r.Context().Value(middleware.ContextUserIdKey).(uint)
		if !ok || userID == 0 {
			zap.L().Error("Error get user_id from context",
				zap.String("op", op))
			res.JSON(w, "unauthorized: user_id not found", http.StatusUnauthorized)
			return
		}

		payloadReq, err := req.JSON[createReq](&w, r)
		if err != nil {
			return
		}

		link := NewLink(userID, payloadReq.URL)
		for {
			existedLink, _ := h.Repository.GetByHash(link.Hash)
			if existedLink == nil {
				break
			}
			link.setHash()
		}

		createdLink, err := h.Repository.Create(link)
		if err != nil {
			zap.L().Error("Error create link",
				zap.String("op", op),
				zap.Error(err))
			res.JSON(w, "error create link", http.StatusInternalServerError)
			return
		}

		fmt.Println("4")

		payloadRes := &createRes{
			Hash: createdLink.Hash,
			URL:  createdLink.Url,
		}

		res.JSON(w, payloadRes, http.StatusCreated)
	}
}

func (h *Handler) goTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "link.handler.goTo()"

		hash := r.PathValue("hash")

		link, err := h.Repository.GetByHash(hash)
		if err != nil {
			zap.L().Error("Error get link by hash",
				zap.String("op", op),
				zap.String("hash", hash),
				zap.Error(err))
		}

		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
	}
}
