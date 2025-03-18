package main

import (
	"URLShorter/configs"
	"URLShorter/internal/link"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func app() http.Handler {
	router := chi.NewRouter()

	link.NewHandler(router)

	return router
}

func main() {
	serverConfig := configs.LoadServerConfig()

	app := app()
	server := &http.Server{
		Addr:    serverConfig.Ip + serverConfig.Port,
		Handler: app,
	}

	fmt.Printf("Server is listening on adress %s:%s\n", serverConfig.Ip, serverConfig.Port)
	err := server.ListenAndServe()
	if err != nil {
		panic("Could not start server: " + err.Error())
	}
}
