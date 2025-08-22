package app

import (
	grpcapp "github.com/wnikx/sso/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GrpcServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// TODO: инициализировахить хранилище
	// TODO: init auth service
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GrpcServer: grpcApp,
	}
}
