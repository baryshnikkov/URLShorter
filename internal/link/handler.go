package link

import (
	"URLShorter/pkg/req"
	"URLShorter/pkg/res"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type handler struct {
}

func NewHandler(router *chi.Mux) {
	handler := &handler{}

	router.Post("/link", handler.create())
	router.Get("/{hash}", handler.goTo())
}

func (h *handler) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payloadReq, err := req.JSON[createReq](w, r)
		if err != nil {
			return
		}

		payloadRes := &createRes{
			URL:  payloadReq.URL,
			Hash: payloadReq.URL,
		}
		res.JSON(w, payloadRes, http.StatusCreated)
	}
}

func (h *handler) goTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")

		fmt.Println(hash)

		http.Redirect(w, r, "https://ya.ru", http.StatusTemporaryRedirect)
	}
}
