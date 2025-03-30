package main

import (
	"URLShorter/configs"
	"URLShorter/internal/link"
	"URLShorter/pkg/database"
	"URLShorter/pkg/logger"
	"URLShorter/pkg/middleware"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

func app() http.Handler {
	appConfig := configs.LoadAppConfig()
	db := database.New(appConfig)
	_ = db

	loggerServeHTTP := logger.New("./logs/shorter/serveHTTP.log")
	router := chi.NewRouter()

	router.Use(middleware.Logger(loggerServeHTTP))
	router.Use(middleware.Gzip)

	link.NewHandler(router)

	return router
}

func main() {
	globalLogger := logger.New("./logs/shorter/errors.log")
	zap.ReplaceGlobals(globalLogger)
	serverConfig := configs.LoadServerConfig()

	app := app()
	server := &http.Server{
		Addr:    serverConfig.Ip + serverConfig.Port,
		Handler: app,
	}

	fmt.Printf("Server is listening on adress %s%s\n", serverConfig.Ip, serverConfig.Port)
	err := server.ListenAndServe()
	if err != nil {
		zap.L().Fatal("Could not start server", zap.Error(err))
	}
}
