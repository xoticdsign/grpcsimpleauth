package app

import (
	"log/slog"
	"time"

	"sso/sso/internal/app/grpcapp"
	"sso/sso/internal/services/auth"
	"sso/sso/internal/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, gRPCPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	gRPCApp := grpcapp.New(log, authService, gRPCPort)

	return &App{
		GRPCServer: gRPCApp,
	}
}
