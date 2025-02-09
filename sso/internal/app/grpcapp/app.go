package grpcapp

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	authgrpc "sso/sso/internal/grpc/auth"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, authService authgrpc.Auth, port int) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	const op = "internal.app.grpcapp.Run()"

	a.log.Info(
		"starting gRPC server",
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		panic("gRPC server couldn't listen: " + err.Error())
	}

	a.log.Info(
		"gRPC server is running",
		slog.String("op", op),
		slog.String("address", l.Addr().String()),
	)

	err = a.gRPCServer.Serve(l)
	if err != nil {
		panic("gRPC server couldn't start: " + err.Error())
	}
}

func (a *App) Stop() {
	const op = "internal.app.grpcapp.Stop()"

	a.log.Info(
		"gRPC server is shutting down",
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	a.gRPCServer.GracefulStop()
}
