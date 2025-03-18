package main

import (
	"fmt"
	"ky-id-backend/src/config"
	mwLogger "ky-id-backend/src/httpserver/middleware/logger"
	"ky-id-backend/src/logger"
	"ky-id-backend/src/storage/sqlite"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	// Check args length
	if len(os.Args) != 3 {
		log.Fatalf("Usage: go run ./src/main.go <config_file> <env_file>")
	}

	fmt.Println("Starting application...")

	// Reading config
	cfg, err := config.MustReadConfig(os.Args[1])

	if err != nil {
		log.Fatalf("%v", err)
	}

	// Reading .env
	err = config.MustReadEnv(os.Args[2])

	if err != nil {
		log.Fatalf("%v", err)
	}

	// Init logger
	_, file, err := logger.InitLogger(cfg)

	if err != nil {
		log.Fatalf("failed to init logger")
	}
	defer file.Close()
	logger.TestLogger()

	// Init storage
	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		slog.Error("failed to init sqlite3 storage", logger.Err(err))
		os.Exit(1)
	}

	_ = storage

	// TODO: init router and handlers
	// Init router
	router := chi.NewRouter()

	// Install middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(mwLogger.LoggerMW)

	// Install routes
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello!\n"))
	})

	// Run server
	var address string = cfg.HTTPServer.Host + ":" + cfg.Port
	server := http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
		WriteTimeout: cfg.Timeout,
	}

	slog.Info("server started", slog.String("address", address))

	if err := server.ListenAndServe(); err != nil {
		slog.Error("failed to start server")
	}

	slog.Error("server stopped")
}
