package app

import (
	"log/slog"

	grpcapp "github.com/paniccaaa/runner/internal/app/grpc"
	"github.com/paniccaaa/runner/internal/services/runner"

	//"github.com/paniccaaa/sso/internal/services/auth"
	"github.com/paniccaaa/runner/internal/storage/postgres"

	ssoGrpc "github.com/paniccaaa/runner/internal/clients/sso/grpc"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewApp(log *slog.Logger, grpcPort int, dbURI string, sso *ssoGrpc.Client) *App {
	// init db
	storage, err := postgres.NewStorage(dbURI)
	if err != nil {
		panic(err)
	}

	// init auth service
	runService := runner.NewRunnerService(log, storage, sso)

	// init grpc
	grpcApp := grpcapp.NewApp(log, grpcPort, runService, storage)

	return &App{
		GRPCServer: grpcApp,
	}
}
