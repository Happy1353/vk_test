package main

import (
	"log/slog"
	"net/http"
	"os"

	"vk/docs"
	"vk/internal/server"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	postgres "vk/internal/storage"

	"github.com/gorilla/mux"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title Your API's Title
// @version 1.0
// @description Your API's Description
// @termsOfService https://example.com/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url https://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /api/v1
func main() {
	log := setupLogger(envLocal)
	log.Debug("Init logger")

	_, err := postgres.New()
	if err != nil {
		log.Error("failed to init storage")
		os.Exit(1)
	}
	log.Info("Init database")

	router := server.SetupRouter()

	InitialSwagger()

	http.ListenAndServe(":8080", router)
}

func InitialSwagger() {
	docs.SwaggerInfo.Schemes = []string{"https", "http"}

	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	go http.ListenAndServe(":8081", r) // Assuming :8081 as Swagger port
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default: // If env config is invalid, set prod settings by default due to security
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
