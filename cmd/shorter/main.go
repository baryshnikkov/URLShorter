package main

import (
	"URLShorter/internal/link"
	"fmt"
	"net/http"
)

const PORT = ":8080"

func app() http.Handler {
	router := http.NewServeMux()

	link.NewHandler(router)

	return router
}

func main() {
	app := app()
	server := &http.Server{
		Addr:    PORT,
		Handler: app,
	}

	fmt.Printf("Server is listening on port %s\n", PORT)
	err := server.ListenAndServe()
	if err != nil {
		panic("Could not start server: " + err.Error())
	}
}
