package app

import (
	"context"
	"log/slog"

	grpcapp "github.com/paniccaaa/runner/internal/app/grpc"
	"github.com/paniccaaa/runner/internal/services/runner"
	"github.com/redis/go-redis/v9"

	//"github.com/paniccaaa/sso/internal/services/auth"
	"github.com/paniccaaa/runner/internal/storage/postgres"

	ssoGrpc "github.com/paniccaaa/runner/internal/clients/sso/grpc"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewApp(log *slog.Logger, grpcPort int, sso *ssoGrpc.Client, dbURI, redisAddr string) *App {
	// init db
	storage, err := postgres.NewStorage(dbURI)
	if err != nil {
		panic(err)
	}

	// init redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	_, err = rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Error("redis connection was refused", slog.String("err", err.Error()))
	}

	// init auth service
	runService := runner.NewRunnerService(log, storage, sso, rdb)

	// init grpc
	grpcApp := grpcapp.NewApp(log, grpcPort, runService, storage, rdb)

	return &App{
		GRPCServer: grpcApp,
	}
}
