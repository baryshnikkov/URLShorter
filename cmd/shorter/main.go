package main

import (
	"URLShorter/configs"
	"URLShorter/internal/auth"
	"URLShorter/internal/link"
	"URLShorter/internal/ping"
	"URLShorter/internal/user"
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

	loggerServeHTTP := logger.New("./logs/shorter/serveHTTP.log")
	router := chi.NewRouter()

	router.Use(middleware.Logger(loggerServeHTTP))
	router.Use(middleware.Gzip)

	userRepository := user.NewRepository(db)

	authService := auth.NewService(userRepository)

	link.NewHandler(router)
	ping.NewHandler(router, &ping.HandlerDeps{Db: db})
	auth.NewHandler(router, &auth.HandlerDeps{Service: authService})

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
