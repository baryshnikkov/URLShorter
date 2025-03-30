package ping

import (
	"URLShorter/pkg/database"
	"URLShorter/pkg/res"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type handler struct {
	db *database.Db
}

type HandlerDeps struct {
	Db *database.Db
}

func NewHandler(router *chi.Mux, deps *HandlerDeps) {
	handler := &handler{
		db: deps.Db,
	}

	router.Get("/ping/db", handler.pingDb())
}

func (h *handler) pingDb() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sqlDB, _ := h.db.SqlDB()
		err := sqlDB.Ping()
		if err != nil {
			zap.L().Error("ping db failed", zap.Error(err))
			res.JSON(w, "Failed to connect to the database", http.StatusInternalServerError)
			return
		}

		res.JSON(w, "Database is successfully available", http.StatusOK)
	}
}
