package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/paniccaaa/runner/internal/app"
	"github.com/paniccaaa/runner/internal/clients/gateway"
	"github.com/paniccaaa/runner/internal/clients/sso/grpc"
	"github.com/paniccaaa/runner/internal/config"
)

func main() {
	cfg := config.MustLoad()

	// init logger
	log := setupLogger(cfg.Env)
	log.Info("starting app", slog.String("env", cfg.Env))

	// init app
	app := app.NewApp(log, cfg.GRPC.Port, cfg.DbURI)

	//init sso-client
	ssoClient, err := grpc.New(
		context.Background(),
		log,
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)
	if err != nil {
		log.Error("failed to init sso client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	isAdmin, err := ssoClient.IsAdmin(context.Background(), 2)
	if err != nil {
		log.Error("is admin error", slog.String("error", err.Error()))
	}

	log.Info("user with user_id=2", slog.Bool("isAdmin", isAdmin))

	// start grpc-server
	go app.GRPCServer.MustRun()

	// start grpc-gateway
	gwServer := gateway.MustRunGRPCGateway(cfg.GRPC.Port, log)
	go gateway.StartGateway(gwServer, log)

	// gracefull shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GRPCServer.Stop()
	gateway.StopGateway(gwServer, log)

	log.Info("App stopped")
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default: //if env config is invalid, set prod settings by default due to security
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
