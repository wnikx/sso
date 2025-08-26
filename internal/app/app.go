package app

import (
	grpcapp "github.com/wnikx/sso/internal/app/grpc"
	"github.com/wnikx/sso/internal/services/auth"
	"github.com/wnikx/sso/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GrpcServer *grpcapp.App
	Auth       *auth.Auth
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}
	authService := auth.New(log, storage, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GrpcServer: grpcApp,
		Auth:       authService,
	}
}
