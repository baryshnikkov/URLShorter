package main

import (
	"URLShorter/internal/link"
	"fmt"
	"net/http"
)

const port = ":8080"

func app() http.Handler {
	router := http.NewServeMux()

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
