package logger

import (
	"ky-id-backend/src/config"
	ph "ky-id-backend/src/logger/handlers/prettyhandler"
	"log/slog"
	"os"
	"time"
)

func InitLogger(cfg *config.Config) (*slog.Logger, *os.File, error) {

	var handler slog.Handler
	var file *os.File

	switch cfg.Enviroment {
	case "local":
		opts := ph.PrettyHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
		handler = ph.NewPrettyHandler(os.Stdout, &opts)
	case "dev":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	case "prod":
		filename := cfg.LogPath + "/logs-" + time.Now().Format("20060102-101010") + ".log"
		file, err := os.Create(filename)

		if err != nil {
			return nil, nil, err
		}

		handler = slog.NewJSONHandler(file, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger, file, nil
}

func TestLogger() {
	slog.Debug(
		"executing database query",
		slog.String("query", "SELECT * FROM users"),
	)
	slog.Info("image upload successful", slog.String("image_id", "39ud88"))
	slog.Warn(
		"storage is 90% full",
		slog.String("available_space", "900.1 MB"),
	)
	slog.Error(
		"An error occurred while processing the request",
		slog.String("url", "https://example.com"),
	)
}
