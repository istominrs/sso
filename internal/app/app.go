package app

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	"sso/internal/repository/postgres"
	"sso/internal/services/auth"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	postgresDSN string,
	tokenTTL time.Duration,
) *App {
	repository, err := postgres.New(postgresDSN)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, repository, repository, repository, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
