package link

import (
	"URLShorter/pkg/req"
	"URLShorter/pkg/res"
	"fmt"
	"net/http"
)

type Handler struct {
}

func NewHandler(router *http.ServeMux) {
	handler := &Handler{}

	router.HandleFunc("POST /link", handler.Create())
	router.HandleFunc("GET /{hash}", handler.GoTo())
}

func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payloadReq, err := req.Payload[CreateReq](w, r)
		if err != nil {
			return
		}

		payloadRes := &CreateRes{
			URL:  payloadReq.URL,
			Hash: payloadReq.URL,
		}
		res.JSON(w, payloadRes, http.StatusCreated)
	}
}

func (h *Handler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")

		fmt.Println(hash)

		http.Redirect(w, r, "https://ya.ru", http.StatusTemporaryRedirect)
	}
}
