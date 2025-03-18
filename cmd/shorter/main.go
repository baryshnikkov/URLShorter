package main

import (
	"URLShorter/internal/link"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const port = ":8080"

func app() http.Handler {
	router := chi.NewRouter()

	link.NewHandler(router)

	return router
}

func main() {
	app := app()
	server := &http.Server{
		Addr:    port,
		Handler: app,
	}

	fmt.Printf("Server is listening on port %s\n", port)
	err := server.ListenAndServe()
	if err != nil {
		panic("Could not start server: " + err.Error())
	}
}
