package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"sso/sso/internal/app"
	"sso/sso/internal/config"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return log
}

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting application",
		slog.Any("config", cfg),
	)

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	application.GRPCServer.Stop()

	log.Info(
		"application stopped",
		slog.String("signal", sign.String()),
	)
}
